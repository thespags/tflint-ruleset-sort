package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CountRule makes sure that `count` attribute is always on top.
type CountRule struct {
	tflint.DefaultRule
}

// NewCountRule creates a new CountRule.
func NewCountRule() *CountRule {
	return &CountRule{}
}

// Name returns the name of the rule.
func (*CountRule) Name() string {
	return project.RuleName("count")
}

// Enabled returns whether the rule is enabled by default.
func (*CountRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*CountRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *CountRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `count` clause is placed on the top of the resource
// definition.
func (r *CountRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "resource" && block.Type != "data" {
			return nil
		}

		count, exists := block.Body.Attributes["count"]
		if !exists {
			return nil
		}

		first := node.FirstNodeFrom(block.Body)
		if first.AsAttribute() == count {
			return nil
		}

		return runner.EmitIssueWithFix(
			r,
			"`count` must be the top-most attribute",
			count.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeBefore(fixer, src, count.SrcRange, first.Range())
			},
		)
	})
}
