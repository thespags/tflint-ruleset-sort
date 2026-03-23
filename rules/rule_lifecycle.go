package rules

import (
	"fmt"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// LifecycleRule makes sure that `lifecycle` clause is the last one before
// `depends_on`.
type LifecycleRule struct {
	tflint.DefaultRule
}

// NewLifecycleRule creates a new LifecycleRule.
func NewLifecycleRule() *LifecycleRule {
	return &LifecycleRule{}
}

// Name returns the name of the rule.
func (*LifecycleRule) Name() string {
	return project.RuleName("lifecycle")
}

// Enabled returns whether the rule is enabled by default.
func (*LifecycleRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*LifecycleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *LifecycleRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the `lifecycle` clause is placed at the end of the
// resource definition (but before `depends_on`).
func (r *LifecycleRule) Check(runner tflint.Runner) error {
	return visit.Blocks(runner, func(block *hclsyntax.Block, src []byte) error {
		if block.Type != "resource" {
			return nil
		}

		lifecycle, err := findBlock(r, runner, block, "lifecycle")
		if lifecycle == nil {
			return err
		}

		nodes := node.OrderedInspectableNodesFrom(block.Body)

		n := nodes[len(nodes)-1]
		if n.IsAttribute() && n.Name() == "depends_on" {
			n = nodes[len(nodes)-2]
		}

		if n.AsBlock() == lifecycle {
			return nil
		}

		return runner.EmitIssueWithFix(
			r,
			"`lifecycle` block must be at the end of the definition (but before `depends_on`)",
			lifecycle.Body.SrcRange,
			func(fixer tflint.Fixer) error {
				// Place lifecycle before depends_on if present, or after the last node
				target := nodes[len(nodes)-1]
				if target.IsAttribute() && target.Name() == "depends_on" {
					return moveNodeBefore(fixer, src, lifecycle.Range(), target.Range())
				}

				return moveNodeAfter(fixer, src, lifecycle.Range(), target.Range())
			},
		)
	})
}

func findBlock(rule tflint.Rule, runner tflint.Runner, block *hclsyntax.Block, name string) (*hclsyntax.Block, error) {
	var value *hclsyntax.Block

	for _, block := range block.Body.Blocks {
		if block.Type == name {
			if value != nil {
				return nil, runner.EmitIssue(
					rule,
					fmt.Sprintf("more than 1 `%s` block found", name),
					block.TypeRange,
				)
			}

			value = block
		}
	}

	return value, nil
}
