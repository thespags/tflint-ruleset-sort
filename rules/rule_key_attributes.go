package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/custom"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// KeyAttributesRule makes sure that key-attributes (those that uniquely
// identify the resource) are put on top of the resource definition.
type KeyAttributesRule struct {
	tflint.DefaultRule
}

// NewKeyAttributesRule creates a new KeyAttributesRule.
func NewKeyAttributesRule() *KeyAttributesRule {
	return &KeyAttributesRule{}
}

// Name returns the name of the rule.
func (*KeyAttributesRule) Name() string {
	return project.RuleName("key_attributes")
}

// Enabled returns whether the rule is enabled by default.
func (*KeyAttributesRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*KeyAttributesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *KeyAttributesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the key-attributes (those that uniquely identify the
// resource) are put on top of the resource definition.
func (r *KeyAttributesRule) Check(rr tflint.Runner) error {
	switch runner := rr.(type) {
	case *custom.Runner:
		return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
			if block.Type != "resource" && block.Type != "data" && block.Type != "module" {
				return nil
			}

			cfg := runner.Lookup(block)
			if cfg == nil {
				return nil
			}

			// For module blocks, `source` precedes key attributes and must
			// not be sorted between them; isResource=false skips it.
			isResource := block.Type != "module"

			return r.checkKeyAttributes(runner, block.Body, src, cfg, isResource)
		})

	default:
		return nil
	}
}

// checkKeyAttributes verifies key attribute ordering within a body.
// isResource controls whether `source` is skipped (false for modules where
// `source` precedes key attributes).
func (r *KeyAttributesRule) checkKeyAttributes(
	runner *custom.Runner,
	body *hclsyntax.Body,
	src []byte,
	resource *custom.Resource,
	isResource bool,
) error {
	knownKeyAttributes := resource.KeyAttributes

	// Navigate into key blocks if needed.
	for _, keyBlockName := range resource.KeyBlocks {
		for _, block := range body.Blocks {
			if block.Type != keyBlockName {
				continue
			}

			body = block.Body
		}
	}

	// Extract key attributes that are present in the definition.
	kaList := make([]*hclsyntax.Attribute, 0, len(knownKeyAttributes))
	kaSet := map[string]struct{}{}

	for _, attrName := range knownKeyAttributes {
		attr, exists := body.Attributes[attrName]
		if !exists {
			continue
		}

		kaList = append(kaList, attr)
		kaSet[attrName] = struct{}{}
	}

	pos := 0

	for _, n := range node.OrderedInspectableNodesFrom(body) {
		// Skip special leading attributes.
		if n.IsAttribute() {
			name := n.Name()
			if name == "for_each" || name == "count" {
				continue
			}

			if !isResource && name == "source" {
				continue
			}
		}

		if pos == len(kaList) {
			break
		}

		attribute := kaList[pos]
		if n.IsAttribute() && n.Name() != attribute.Name {
			nRange := n.Range()
			if _, isKey := kaSet[n.Name()]; isKey {
				return runner.EmitIssueWithFix(
					r,
					fmt.Sprintf(
						"higher-priority key-attribute `%s` should be defined before `%s`",
						attribute.Name,
						n.Name(),
					),
					attribute.SrcRange,
					func(fixer tflint.Fixer) error {
						return moveNodeBefore(fixer, src, attribute.SrcRange, nRange)
					},
				)
			}

			return runner.EmitIssueWithFix(
				r,
				fmt.Sprintf(
					"key-attribute `%s` should be defined before non-key `%s`",
					attribute.Name,
					n.Name(),
				),
				attribute.SrcRange,
				func(fixer tflint.Fixer) error {
					return moveNodeBefore(fixer, src, attribute.SrcRange, nRange)
				},
			)
		}

		pos++
	}

	return nil
}
