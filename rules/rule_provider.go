package rules

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// ProviderRule makes sure that `provider` is placed after for_each/count but
// before all other attributes and blocks.
type ProviderRule struct {
	tflint.DefaultRule
}

// NewProviderRule creates a new ProviderRule.
func NewProviderRule() *ProviderRule {
	return &ProviderRule{}
}

// Name returns the name of the rule.
func (*ProviderRule) Name() string {
	return project.RuleName("provider")
}

// Enabled returns whether the rule is enabled by default.
func (*ProviderRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*ProviderRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *ProviderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `provider` attribute is placed after for_each/count
// but before everything else.
func (r *ProviderRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "resource" && block.Type != "data" {
			return nil
		}

		provider, exists := block.Body.Attributes["provider"]
		if !exists {
			return nil
		}

		nodes := node.OrderedInspectableNodesFrom(block.Body)

		// Skip leading for_each and count — provider comes after them.
		for len(nodes) > 0 {
			a := nodes[0].AsAttribute()
			if a == nil || (a.Name != "for_each" && a.Name != "count") {
				break
			}

			nodes = nodes[1:]
		}

		if len(nodes) > 0 && nodes[0].AsAttribute() == provider {
			return nil
		}

		// provider is not at the expected position — find the target to move before.
		target := nodes[0]

		return runner.EmitIssueWithFix(
			r,
			"`provider` must follow `for_each`/`count` (or be the top-most attribute)",
			provider.SrcRange,
			func(fixer tflint.Fixer) error {
				return moveNodeBefore(fixer, src, provider.SrcRange, target.Range())
			},
		)
	})
}
