package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
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

// Check verifies whether the `source` clause is placed on the top of the module
// definition.
func (r *SourceRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "module" {
			return nil
		}

		source, exists := block.Body.Attributes["source"]
		if !exists {
			return nil
		}

		first := node.FirstNodeFrom(block.Body)
		if first.AsAttribute() == source {
			return nil
		}

		return runner.EmitIssueWithFix(
			r,
			"`source` must be the top-most attribute",
			source.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeBefore(fixer, src, source.SrcRange, first.Range())
			},
		)
	})
}
