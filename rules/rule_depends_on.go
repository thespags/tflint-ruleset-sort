package rules

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// DependsOnRule makes sure that `depends_on` clause is always the last.
type DependsOnRule struct {
	tflint.DefaultRule
}

// NewDependsOnRule creates a new DependsOnRule.
func NewDependsOnRule() *DependsOnRule {
	return &DependsOnRule{}
}

// Name returns the name of the rule.
func (*DependsOnRule) Name() string {
	return project.RuleName("depends_on")
}

// Enabled returns whether the rule is enabled by default.
func (*DependsOnRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*DependsOnRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *DependsOnRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `depends_on` clause is placed at the end of the
// resource definition.
func (r *DependsOnRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "resource" && block.Type != "data" {
			return nil
		}

		dependsOn, exists := block.Body.Attributes["depends_on"]
		if !exists {
			return nil
		}

		last := node.LastNodeFrom(block.Body)
		if last.AsAttribute() == dependsOn {
			return nil
		}

		return runner.EmitIssueWithFix(
			r,
			"`depends_on` clause must be the last one in the definition",
			dependsOn.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeAfter(fixer, src, dependsOn.SrcRange, last.Range())
			},
		)
	})
}
