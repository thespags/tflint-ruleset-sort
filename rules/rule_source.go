package rules

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// SourceRule makes sure that `source` meta-attribute is always on top.
type SourceRule struct {
	tflint.DefaultRule
}

// NewSourceRule creates a new SourceRule.
func NewSourceRule() *SourceRule {
	return &SourceRule{}
}

// Name returns the name of the rule.
func (*SourceRule) Name() string {
	return project.RuleName("source")
}

// Enabled returns whether the rule is enabled by default.
func (*SourceRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*SourceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `source` clause is placed after for_each/count
// but before everything else in the module definition.
func (r *SourceRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "module" {
			return nil
		}

		source, exists := block.Body.Attributes["source"]
		if !exists {
			return nil
		}

		nodes := node.OrderedInspectableNodesFrom(block.Body)

		// Skip leading for_each and count — source comes after them.
		for len(nodes) > 0 {
			a := nodes[0].AsAttribute()
			if a == nil || (a.Name != "for_each" && a.Name != "count") {
				break
			}

			nodes = nodes[1:]
		}

		if len(nodes) > 0 && nodes[0].AsAttribute() == source {
			return nil
		}

		// source is not at the expected position — find the target to move before.
		target := nodes[0]

		return runner.EmitIssueWithFix(
			r,
			"`source` must follow `for_each`/`count` (or be the top-most attribute)",
			source.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeBefore(fixer, src, source.SrcRange, target.Range())
			},
		)
	})
}
