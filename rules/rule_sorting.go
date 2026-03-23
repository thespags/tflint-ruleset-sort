package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

type optSorting uint

const (
	dontSortMultiliners optSorting = 1 << iota
)

func (o optSorting) dontSortMultiliners() bool {
	return o&dontSortMultiliners != 0
}

// SortingRule makes sure that all attributes and blocks are properly sorted.
type SortingRule struct {
	tflint.DefaultRule
}

// NewSortingRule creates a new SortingRule.
func NewSortingRule() *SortingRule {
	return &SortingRule{}
}

// Name returns the name of the rule.
func (*SortingRule) Name() string {
	return project.RuleName("sorting")
}

// Enabled returns whether the rule is enabled by default.
func (*SortingRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*SortingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SortingRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether all attributes and blocks are properly sorted.
func (r *SortingRule) Check(rr tflint.Runner) error {
	switch runner := rr.(type) {
	case *custom.Runner:
		return visit.Files(runner, func(b *hclsyntax.Body, src []byte) error {
			return r.stepInto(runner, src, 0, node.OrderedInspectableNodesFrom(b))
		})

	default:
		return nil
	}
}

func (r *SortingRule) stepInto(
	runner *custom.Runner,
	src []byte,
	level int,
	nodes []node.InspectableNode,
) error {
	for _, n := range nodes {
		if a := n.AsAttribute(); a != nil {
			if err := r.checkExpression(runner, level+1, src, a.Expr); err != nil {
				return err
			}
		} else if b := n.AsBlock(); b != nil {
			if err := r.checkBlock(runner, src, level+1, b); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *SortingRule) checkNodes(
	runner *custom.Runner,
	src []byte,
	level int,
	nodes []node.InspectableNode,
) error {
	text := "single-line nodes `(%s)` should precede multi-lines `(%s)`"
	if fixed, err := r.reorder(runner, src, nodes, singleLinesFirst, text); fixed || err != nil {
		return err
	}

	text = "attributes `%s` should precede blocks `%s`"
	if fixed, err := r.reorder(runner, src, nodes, attributesFirst, text); fixed || err != nil {
		return err
	}

	if fixed, err := r.sortAlphabetically(runner, src, nodes); fixed || err != nil {
		return err
	}

	return r.stepInto(runner, src, level, nodes)
}

func (r *SortingRule) sortAlphabetically(
	runner *custom.Runner,
	src []byte,
	nodes []node.InspectableNode,
) (bool, error) {
	sorted := slices.Clone(nodes)
	if len(sorted) < 2 {
		return false, nil
	}

	start := 0

	i := 1
	for ; i < len(nodes); i++ {
		left := nodes[i-1]
		right := nodes[i]

		var endGroup bool

		if left.Lines() > 1 && right.Lines() > 1 {
			// For multi-line nodes, we also split for comments... (legacy).
			var err error

			endGroup, err = hasCommentsInBetween(src, left, right)
			if err != nil {
				return false, err
			}
		} else {
			// You can split single lines by new lines into separate groups.
			endGroup = right.Range().Start.Line-left.Range().End.Line > 1
		}
		// We sort sections of the same lines, the same kind, or special cases above.
		// If those don't match, then we sort the group of nodes and find the next group.
		if endGroup || left.Kind() != right.Kind() {
			slices.SortStableFunc(sorted[start:i], nameOrdered)
			start = i

			continue
		}
	}

	slices.SortStableFunc(sorted[start:i], nameOrdered)

	// If nothing changed, we don't emit an issue.
	if slices.EqualFunc(nodes, sorted, func(a, b node.InspectableNode) bool {
		return a == b
	}) {
		return false, nil
	}

	return true, runner.EmitIssueWithFix(
		r,
		fmt.Sprintf("`(%s)` should be reordered `(%s)` (alphabetical sorting)", toNames(nodes), toNames(sorted)),
		sortedRange(nodes, sorted),
		func(fixer tflint.Fixer) error {
			return reorderNodes(fixer, src, nodes, sorted)
		},
	)
}

func (r *SortingRule) reorder(
	runner *custom.Runner,
	src []byte,
	nodes []node.InspectableNode,
	cmp func(a, b node.InspectableNode) int,
	text string,
) (bool, error) {
	lasts := make([]node.InspectableNode, 0)
	firsts := make([]node.InspectableNode, 0)

	for i := 1; i < len(nodes); i++ {
		left := nodes[i-1]
		right := nodes[i]

		if cmp(left, right) == 1 {
			lasts = append(lasts, left)
			firsts = append(firsts, right)
		}
	}

	if len(lasts) == 0 {
		return false, nil
	}

	span := hcl.RangeBetween(lasts[0].Range(), firsts[len(firsts)-1].Range())

	sorted := slices.Clone(nodes)
	slices.SortStableFunc(sorted, func(a, b node.InspectableNode) int {
		if a.Lines() > 1 && b.Lines() == 1 {
			return 1
		}

		return -1
	})

	return true, runner.EmitIssueWithFix(
		r,
		fmt.Sprintf(text, toNames(firsts), toNames(lasts)),
		span,
		func(fixer tflint.Fixer) error {
			return reorderNodes(fixer, src, nodes, sorted)
		},
	)
}

func (r *SortingRule) checkBlock(
	runner *custom.Runner,
	src []byte,
	level int,
	block *hclsyntax.Block,
) error {
	nodes := node.OrderedInspectableNodesFrom(block.Body)
	if len(nodes) == 0 {
		return nil
	}

	if level == 1 {
		// Top-level `locals` block: just check the expressions and exit
		if block.Type == "locals" {
			for _, _n := range nodes {
				if _a := _n.AsAttribute(); _a != nil {
					if err := r.checkExpression(runner, level+1, src, _a.Expr, dontSortMultiliners); err != nil {
						return err
					}
				}
			}

			return nil
		}

		// Top-level `resource` or `data` block
		if block.Type == "resource" || block.Type == "data" {
			var err error

			nodes, err = r.preprocessResourceOrData(runner, src, level, block.Labels[0], nodes)
			if err != nil {
				return err
			}
		}

		// Top-level `module` block
		if block.Type == "module" {
			nodes = r.preprocessModule(runner, block.Body, nodes)
		}
	}

	return r.checkNodes(runner, src, level+1, nodes)
}

func (r *SortingRule) preprocessResourceOrData(
	runner *custom.Runner,
	src []byte,
	level int,
	kind string,
	nodes []node.InspectableNode,
) ([]node.InspectableNode, error) {
	// Drop leading `for_each`
	if disabled, exists := runner.Disabled["for_each"]; !exists || !disabled {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "for_each" {
			nodes = nodes[1:]
		}

		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop leading `count`
	if disabled, exists := runner.Disabled["count"]; !exists || !disabled {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "count" {
			nodes = nodes[1:]
		}

		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop trailing `depends_on`
	if disabled, exists := runner.Disabled["depends_on"]; !exists || disabled {
		if a := nodes[len(nodes)-1].AsAttribute(); a != nil && a.Name == "depends_on" {
			nodes = nodes[:len(nodes)-1]
		}

		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop trailing `lifecycle`
	if disabled, exists := runner.Disabled["lifecycle"]; !exists || !disabled {
		if b := nodes[len(nodes)-1].AsBlock(); b != nil && b.Type == "lifecycle" {
			nodes = nodes[:len(nodes)-1]
		}

		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop key-attributes
	if disabled, exists := runner.Disabled["key_attributes"]; !exists || !disabled {
		// Are there key-attributes?
		resource, ok := runner.Resources[kind]
		if !ok {
			return nodes, nil
		}

		if len(resource.KeyAttributes) == 0 {
			return nodes, nil
		}

		// Drop key-attributes
		if level > len(resource.KeyBlocks) {
			for _, key := range resource.KeyAttributes {
				if len(nodes) == 0 {
					break
				}

				if a := nodes[0].AsAttribute(); a != nil {
					if a.Name == key {
						nodes = nodes[1:] // Drop
					}
				}
			}

			return nodes, nil
		}

		// Drop the key block but inspect its contents
		// KeyBlocks are the mapping of prefix, so manifest.metadata -> ["manifest", "metadata"]
		// And, we're only allowed one set of KeyBlocks.
		if kb := nodes[0].AsBlock(); kb != nil && kb.Type == resource.KeyBlocks[level-1] {
			nodes = nodes[1:] // Drop
			knodes := node.OrderedInspectableNodesFrom(kb.Body)
			// Drop leading key attributes inside the key block
			if len(resource.KeyBlocks) == level {
				for _, key := range resource.KeyAttributes {
					if len(knodes) == 0 {
						break
					}

					if a := knodes[0].AsAttribute(); a != nil {
						if a.Name == key {
							knodes = knodes[1:] // Drop
						}
					}
				}
			}
			// Step inside
			if err := r.checkNodes(runner, src, level+1, knodes); err != nil {
				return nil, err
			}
		}
	}

	return nodes, nil
}

func (*SortingRule) preprocessModule(
	runner *custom.Runner,
	body *hclsyntax.Body,
	nodes []node.InspectableNode,
) []node.InspectableNode {
	// Drop leading `source`
	if disabled, exists := runner.Disabled["source"]; !exists || !disabled {
		if len(nodes) > 0 {
			if a := nodes[0].AsAttribute(); a != nil && a.Name == "source" {
				nodes = nodes[1:]
			}
		}

		if len(nodes) == 0 {
			return nodes
		}
	}

	// Drop key-attributes for modules (keyed by source value)
	if disabled, exists := runner.Disabled["key_attributes"]; !exists || !disabled {
		source, exists := body.Attributes["source"]
		if exists {
			val, diags := source.Expr.Value(nil)
			if !diags.HasErrors() && val.Type() == cty.String {
				if module, ok := runner.Modules[val.AsString()]; ok {
					for _, key := range module.KeyAttributes {
						if len(nodes) == 0 {
							break
						}

						if a := nodes[0].AsAttribute(); a != nil && a.Name == key {
							nodes = nodes[1:]
						}
					}
				}
			}
		}
	}

	return nodes
}

func (r *SortingRule) checkExpression(
	runner *custom.Runner,
	level int,
	src []byte,
	expression hclsyntax.Expression,
	opts ...optSorting,
) error {
	var opt optSorting
	for _, o := range opts {
		opt |= o
	}

	switch expr := expression.(type) {
	case *hclsyntax.ForExpr:
		return r.checkForExpr(runner, level, src, expr, opt)
	case *hclsyntax.FunctionCallExpr:
		return r.checkFunctionCallExpr(runner, level, src, expr, opt)
	case *hclsyntax.ObjectConsExpr:
		return r.checkObjectConsExpr(runner, level, src, expr, opt)
	case *hclsyntax.ParenthesesExpr:
		return r.checkParenthesesExpr(runner, level, src, expr, opt)
	case *hclsyntax.TupleConsExpr:
		return r.checkTupleConsExpr(runner, level, src, expr, opt)
	default:
		return nil
	}
}

func (r *SortingRule) checkForExpr(
	runner *custom.Runner,
	level int,
	src []byte,
	x *hclsyntax.ForExpr,
	opt optSorting,
) error {
	// Step inside
	return r.checkExpression(runner, level+1, src, x.ValExpr, opt)
}

func (r *SortingRule) checkFunctionCallExpr(
	runner *custom.Runner,
	level int,
	src []byte,
	exor *hclsyntax.FunctionCallExpr,
	opt optSorting,
) error {
	// Step inside
	for _, e := range exor.Args {
		if err := r.checkExpression(runner, level+1, src, e, opt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SortingRule) checkObjectConsExpr(
	runner *custom.Runner,
	level int,
	src []byte,
	expr *hclsyntax.ObjectConsExpr,
	opt optSorting,
) error {
	for i := 1; i < len(expr.Items); i++ {
		curX := expr.Items[i]
		prevX := expr.Items[i-1]

		curLines := curX.ValueExpr.Range().End.Line - curX.KeyExpr.Range().Start.Line
		prevLines := prevX.ValueExpr.Range().End.Line - prevX.KeyExpr.Range().Start.Line

		curKey := strings.ToLower(node.WrapObjectConsItem(&curX).Name())
		prevKey := strings.ToLower(node.WrapObjectConsItem(&prevX).Name())

		if prevLines > 0 && curLines == 0 {
			prevRange := node.WrapObjectConsItem(&prevX).Range()
			curRange := node.WrapObjectConsItem(&curX).Range()

			if err := runner.EmitIssueWithFix(
				r,
				"single-line key/value pair should be placed before multi-line",
				curX.KeyExpr.Range(),
				func(fixer tflint.Fixer) error {
					return swapAdjacentNodes(fixer, src, prevRange, curRange)
				},
			); err != nil {
				return err
			}
		}

		if prevLines == 0 && curLines == 0 &&
			curX.KeyExpr.Range().Start.Line-prevX.ValueExpr.Range().End.Line == 1 &&
			curKey < prevKey {
			prevRange := node.WrapObjectConsItem(&prevX).Range()
			curRange := node.WrapObjectConsItem(&curX).Range()

			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf(
					"key `%s` is out of order (should not follow alphabetically greater `%s`)",
					curKey,
					prevKey,
				),
				curX.KeyExpr.Range(),
				func(fixer tflint.Fixer) error {
					return swapAdjacentNodes(fixer, src, prevRange, curRange)
				},
			); err != nil {
				return err
			}
		}

		if prevLines > 0 && curLines > 0 && !opt.dontSortMultiliners() && curKey < prevKey {
			prevRange := node.WrapObjectConsItem(&prevX).Range()
			curRange := node.WrapObjectConsItem(&curX).Range()

			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf(
					"key `%s` is out of order (should not follow alphabetically greater `%s`)",
					curKey,
					prevKey,
				),
				curX.KeyExpr.Range(),
				func(fixer tflint.Fixer) error {
					return swapAdjacentNodes(fixer, src, prevRange, curRange)
				},
			); err != nil {
				return err
			}
		}
	}

	// Step inside
	for _, e := range expr.Items {
		if err := r.checkExpression(runner, level+1, src, e.ValueExpr, opt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SortingRule) checkTupleConsExpr(
	runner *custom.Runner,
	level int,
	src []byte,
	x *hclsyntax.TupleConsExpr,
	opt optSorting,
) error {
	// Step inside
	for _, e := range x.Exprs {
		if err := r.checkExpression(runner, level+1, src, e, opt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SortingRule) checkParenthesesExpr(
	runner *custom.Runner,
	level int,
	src []byte,
	x *hclsyntax.ParenthesesExpr,
	opt optSorting,
) error {
	// Step inside
	return r.checkExpression(runner, level+1, src, x.Expression, opt)
}
