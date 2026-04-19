package rules

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// ForEachRule makes sure that `for_each` meta-attribute is always on top.
type ForEachRule struct {
	tflint.DefaultRule
}

// NewForEachRule creates a new ForEachRule.
func NewForEachRule() *ForEachRule {
	return &ForEachRule{}
}

// Name returns the name of the rule.
func (*ForEachRule) Name() string {
	return project.RuleName("for_each")
}

// Enabled returns whether the rule is enabled by default.
func (*ForEachRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*ForEachRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *ForEachRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `for_each` clause is placed on the top of the
// resource, data, or module definition.
func (r *ForEachRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "resource" && block.Type != "data" && block.Type != "module" {
			return nil
		}

		forEach, exists := block.Body.Attributes["for_each"]
		if !exists {
			return nil
		}

		first := node.FirstNodeFrom(block.Body)
		if first.AsAttribute() == forEach {
			return nil
		}

		return runner.EmitIssueWithFix(
			r,
			"`for_each` must be the top-most attribute",
			forEach.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeBefore(fixer, src, forEach.SrcRange, first.Range())
			},
		)
	})
}
