package rules

import (
	"fmt"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// UnknownResourceRule warns if the linter encounters unknown resource.
type UnknownResourceRule struct {
	tflint.DefaultRule

	seenResources map[string]struct{}
	seenModules   map[string]struct{}
}

// NewUnknownResourceRule creates a new UnknownResourceRule.
func NewUnknownResourceRule() *UnknownResourceRule {
	return &UnknownResourceRule{
		seenResources: map[string]struct{}{},
		seenModules:   map[string]struct{}{},
	}
}

// Name returns the name of the rule.
func (*UnknownResourceRule) Name() string {
	return project.RuleName("unknown_resource")
}

// Enabled returns whether the rule is enabled by default.
func (*UnknownResourceRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*UnknownResourceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the reference link for the rule.
func (r *UnknownResourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the key-attributes (those that uniquely identify the
// resource) are put on top of the resource definition.
func (r *UnknownResourceRule) Check(rr tflint.Runner) error {
	switch runner := rr.(type) {
	case *custom.Runner:
		return visit.Blocks(runner, func(block *hclsyntax.Block, _ []byte) error {
			if block.Type == "resource" || block.Type == "data" {
				kind := block.Labels[0]
				if _, kindIsKnown := runner.Resources[kind]; kindIsKnown {
					return nil
				}

				if _, ok := r.seenResources[kind]; !ok {
					r.seenResources[kind] = struct{}{}

					return runner.EmitIssue(
						r,
						fmt.Sprintf("key-attributes for resource type `%s` are not configured", kind),
						block.LabelRanges[0],
					)
				}

				return nil
			}

			if block.Type == "module" {
				kind := block.Labels[0]
				if _, kindIsKnown := runner.Modules[kind]; kindIsKnown {
					return nil
				}

				if _, ok := r.seenModules[kind]; !ok {
					r.seenModules[kind] = struct{}{}

					return runner.EmitIssue(
						r,
						fmt.Sprintf("key-attributes for resource type `%s` are not configured", kind),
						block.LabelRanges[0],
					)
				}

				return nil
			}

			return nil
		})

	default:
		return nil
	}
}
