package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/custom"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

const (
	separateAttribute      = "attribute `%s` must be separated from the rest of the definition by an extra line"
	noSeparateKeyAttribute = "key-attribute `%s` must not be separated from `source` or other key-attributes"
	separateMultiLine      = "multi-line element must be separated from the previous one by an extra line"
	separateSingleLine     = "single-line element must be separated from the preceding multi-line one by an extra line"
)

// SpacingRule makes sure that there is consistent spacing between attributes
// blocks and.
type SpacingRule struct {
	tflint.DefaultRule
}

// NewSpacingRule creates a new SpacingRule.
func NewSpacingRule() *SpacingRule {
	return &SpacingRule{}
}

// Name returns the name of the rule.
func (*SpacingRule) Name() string {
	return project.RuleName("spacing")
}

// Enabled returns whether the rule is enabled by default.
func (*SpacingRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*SpacingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SpacingRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether all attributes and blocks are properly sorted.
func (r *SpacingRule) Check(rr tflint.Runner) error {
	switch runner := rr.(type) {
	case *custom.Runner:
		return visit.Files(runner, func(b *hclsyntax.Body, src []byte) error {
			return r.checkNodes(runner, src, 0, node.OrderedInspectableNodesFrom(b), b)
		})

	default:
		return nil
	}
}

func (r *SpacingRule) checkNodes(
	runner *custom.Runner,
	src []byte,
	level int,
	nodes []node.InspectableNode,
	parent node.Node,
) error {
	if len(nodes) == 0 {
		return r.checkEmptySpace(runner, src, parent)
	}

	if level == 0 {
		if err := r.checkLeadingSpace(runner, src, nodes[0], parent); err != nil {
			return err
		}
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(nodes); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			nodes[i],
			nodes[i-1],
			isolateMultiLiners,
		); err != nil {
			return err
		}
	}

	// Step inside
	for _, n := range nodes {
		if b := n.AsBlock(); b != nil {
			if err := r.checkBlock(runner, src, level+1, b); err != nil {
				return err
			}
		} else if a := n.AsAttribute(); a != nil {
			if err := r.checkExpression(runner, src, level+1, a.Expr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *SpacingRule) checkBlock(
	runner *custom.Runner,
	src []byte,
	level int,
	block *hclsyntax.Block,
) error {
	nodes := node.OrderedInspectableNodesFrom(block.Body)

	if len(nodes) > 0 {
		// Check the leading empty lines
		if err := r.checkLeadingSpace(runner, src, nodes[0], block.Body); err != nil {
			return err
		}
		// Check the trailing empty lines
		if err := r.checkTrailingSpace(runner, src, nodes[len(nodes)-1], block.Body); err != nil {
			return err
		}
	}

	// Check special treatment of `count`, `for_each`, `depends_on`, and `lifecycle`
	if level == 1 && (block.Type == "resource" || block.Type == "data") {
		if len(nodes) > 1 {
			check := func(left, _ node.InspectableNode) (bool, string, hcl.Range) {
				a := left.AsAttribute()

				return a != nil && (a.Name == "count" || a.Name == "for_each"),
					fmt.Sprintf(separateAttribute, left.Name()),
					left.Range()
			}

			emitted, err := r.requireLineBetween(runner, src, check, nodes[0], nodes[1])
			if err != nil {
				return err
			}

			if emitted {
				nodes = nodes[1:]
			}
		}

		if len(nodes) > 1 {
			// Trailing depends_on: must be separated by a blank line from the rest.
			end := len(nodes) - 1
			check := func(_, right node.InspectableNode) (bool, string, hcl.Range) {
				a := right.AsAttribute()

				return a != nil && a.Name == "depends_on",
					fmt.Sprintf(separateAttribute, right.Name()),
					right.Range()
			}

			emitted, err := r.requireLineBetween(runner, src, check, nodes[end-1], nodes[end])
			if err != nil {
				return err
			}

			if emitted {
				nodes = nodes[:end]
			}
		}

		if len(nodes) > 1 {
			// Trailing lifecycle: must be separated by a blank line from the rest.
			end := len(nodes) - 1
			check := func(_, right node.InspectableNode) (bool, string, hcl.Range) {
				innerBlock := right.AsBlock()

				return innerBlock != nil && innerBlock.Type == "lifecycle",
					fmt.Sprintf(separateAttribute, right.Name()),
					right.Range()
			}

			emitted, err := r.requireLineBetween(runner, src, check, nodes[end-1], nodes[end])
			if err != nil {
				return err
			}

			if emitted {
				nodes = nodes[:end]
			}
		}

		// Check key-attribute spacing for resources/data
		if len(nodes) > 1 {
			kind := block.Labels[0]

			var err error

			nodes, err = r.checkKeyAttributeSpacing(runner, src, nodes, keyAttrSet(runner.Resources, kind))
			if len(nodes) == 0 || err != nil {
				return err
			}
		}
	}

	// Check special treatment of `source` and key-attributes for modules
	if level == 1 && len(nodes) > 1 && block.Type == "module" {
		var err error

		nodes, err = r.checkModuleSpacing(runner, block, src, nodes)
		if len(nodes) == 0 || err != nil {
			return err
		}
	}

	// Check inside the nested nodes
	return r.checkNodes(runner, src, level, nodes, block.Body)
}

func (r *SpacingRule) checkExpression(
	runner *custom.Runner,
	src []byte,
	level int,
	expression hclsyntax.Expression,
) error {
	switch expr := expression.(type) {
	case *hclsyntax.ForExpr:
		return r.checkForExpr(runner, src, level, expr)
	case *hclsyntax.FunctionCallExpr:
		return r.checkFunctionCallExpr(runner, src, level, expr)
	case *hclsyntax.ObjectConsExpr:
		return r.checkObjectConsExpr(runner, src, level, expr)
	case *hclsyntax.ParenthesesExpr:
		return r.checkParenthesesExpr(runner, src, expr)
	case *hclsyntax.TupleConsExpr:
		return r.checkTupleConsExpr(runner, src, level, expr)
	default:
		return nil
	}
}

func (r *SpacingRule) checkForExpr(
	runner *custom.Runner,
	src []byte,
	level int,
	x *hclsyntax.ForExpr,
) error {
	// Step inside
	return r.checkExpression(runner, src, level+1, x.ValExpr)
}

func (r *SpacingRule) checkFunctionCallExpr(
	runner *custom.Runner,
	src []byte,
	level int,
	expr *hclsyntax.FunctionCallExpr,
) error {
	// Check empty arguments list
	if len(expr.Args) == 0 {
		return r.checkEmptySpace(runner, src, expr)
	}

	if err := r.checkLeadingSpace(runner, src, expr.Args[0], expr); err != nil {
		return err
	}

	if err := r.checkTrailingSpace(runner, src, expr.Args[len(expr.Args)-1], expr); err != nil {
		return err
	}

	for i := 1; i < len(expr.Args); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			node.WrapExpression(expr.Args[i]),
			node.WrapExpression(expr.Args[i-1]),
			dontIsolateMultiLiners,
		); err != nil {
			return err
		}
	}

	// Step inside
	for _, e := range expr.Args {
		if err := r.checkExpression(runner, src, level+1, e); err != nil {
			return err
		}
	}

	return nil
}

func (r *SpacingRule) checkObjectConsExpr(
	runner *custom.Runner,
	src []byte,
	level int,
	expr *hclsyntax.ObjectConsExpr,
) error {
	// Check empty object
	if len(expr.Items) == 0 {
		return r.checkEmptySpace(runner, src, expr)
	}

	// Check the leading empty lines
	if err := r.checkLeadingSpace(runner, src, expr.Items[0].KeyExpr, expr); err != nil {
		return err
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(expr.Items); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			node.WrapObjectConsItem(&expr.Items[i]),
			node.WrapObjectConsItem(&expr.Items[i-1]),
			isolateMultiLiners,
		); err != nil {
			return err
		}
	}

	// Check the trailing empty lines
	if err := r.checkTrailingSpace(runner, src, expr.Items[len(expr.Items)-1].ValueExpr, expr); err != nil {
		return err
	}

	// Step inside
	for _, item := range expr.Items {
		if err := r.checkExpression(
			runner,
			src,
			level+1,
			item.ValueExpr,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *SpacingRule) checkParenthesesExpr(
	runner *custom.Runner,
	src []byte,
	expr *hclsyntax.ParenthesesExpr,
) error {
	if err := r.checkLeadingSpace(runner, src, expr.Expression, expr); err != nil {
		return err
	}

	return r.checkTrailingSpace(runner, src, expr.Expression, expr)
}

func (r *SpacingRule) checkTupleConsExpr(
	runner *custom.Runner,
	src []byte,
	level int,
	expr *hclsyntax.TupleConsExpr,
) error {
	// Check empty list
	if len(expr.Exprs) == 0 {
		return r.checkEmptySpace(runner, src, expr)
	}

	// Check the leading empty lines
	if err := r.checkLeadingSpace(runner, src, expr.Exprs[0], expr); err != nil {
		return err
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(expr.Exprs); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			node.WrapExpression(expr.Exprs[i]),
			node.WrapExpression(expr.Exprs[i-1]),
			dontIsolateMultiLiners,
		); err != nil {
			return err
		}
	}

	// Check the trailing empty lines
	if err := r.checkTrailingSpace(runner, src, expr.Exprs[len(expr.Exprs)-1], expr); err != nil {
		return err
	}

	// Step inside
	for _, innerExpr := range expr.Exprs {
		if err := r.checkExpression(runner, src, level+1, innerExpr); err != nil {
			return err
		}
	}

	return nil
}

// checkEmptySpace inspects the empty space within the range of the node.
func (r *SpacingRule) checkEmptySpace(
	runner *custom.Runner,
	src []byte,
	n node.Node,
) error {
	lines := n.Range().End.Line - n.Range().Start.Line

	if lines > 1 {
		_src := src[n.Range().Start.Byte:n.Range().End.Byte]
		_rng := hcl.Range{
			Filename: n.Range().Filename,
			Start:    n.Range().Start,
			End:      n.Range().End,
		}

		if err := r.checkLines(runner, src, _src, _rng, func(lines int, _ hcl.Range, _ bool) error {
			if lines > 1 {
				excess := lines - 1

				ndRange := n.Range()
				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf(
						"%d redundant blank %s between the braces",
						excess,
						lineOrLines(excess),
					),
					n.Range(),
					func(fixer tflint.Fixer) error {
						return removeBlankLinesInRegion(fixer, src, ndRange, excess)
					},
				); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// checkLeadingSpace inspects the empty space in front of the first node inside
// the parent's range.
func (r *SpacingRule) checkLeadingSpace(
	runner *custom.Runner,
	src []byte,
	child node.Node,
	parent node.Node,
) error {
	lines := child.Range().Start.Line - parent.Range().Start.Line
	if lines == 0 {
		return nil
	}

	_src := src[parent.Range().Start.Byte:child.Range().Start.Byte]

	_rng := hcl.Range{
		Filename: parent.Range().Filename,
		Start:    parent.Range().Start,
		End:      child.Range().Start,
	}

	return r.checkLines(runner, src, _src, _rng, func(lines int, _ hcl.Range, sawComments bool) error {
		threshold := 0
		if sawComments {
			threshold = 1
		}

		if lines > threshold {
			excess := lines - threshold

			betweenRegion := hcl.Range{
				Filename: parent.Range().Filename,
				Start:    parent.Range().Start,
				End:      child.Range().Start,
			}
			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf(
					"%d redundant blank %s in front",
					excess,
					lineOrLines(excess),
				),
				child.Range(),
				func(fixer tflint.Fixer) error {
					return removeBlankLinesInRegion(fixer, src, betweenRegion, excess)
				},
			); err != nil {
				return err
			}
		}

		return nil
	})
}

type optIsolateMultiLiners int

const (
	isolateMultiLiners optIsolateMultiLiners = iota
	dontIsolateMultiLiners
)

// checkMiddleSpace inspects the contents in-between the nodes.
func (r *SpacingRule) checkMiddleSpace(
	runner *custom.Runner,
	src []byte,
	cur node.InspectableNode,
	prev node.InspectableNode,
	opt optIsolateMultiLiners,
) error {
	lines := cur.Range().Start.Line - prev.Range().End.Line

	if lines > 1 {
		_src := src[prev.Range().End.Byte:cur.Range().Start.Byte]

		_rng := hcl.Range{
			Filename: prev.Range().Filename,
			Start:    prev.Range().End,
			End:      cur.Range().Start,
		}
		if err := r.checkLines(runner, src, _src, _rng, func(lines int, _ hcl.Range, _ bool) error {
			if lines > 1 {
				excess := lines - 1

				betweenRegion := hcl.Range{
					Filename: prev.Range().Filename,
					Start:    prev.Range().End,
					End:      cur.Range().Start,
				}
				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf(
						"%d redundant blank %s in front",
						excess,
						lineOrLines(excess),
					),
					cur.Range(),
					func(fixer tflint.Fixer) error {
						return removeBlankLinesInRegion(fixer, src, betweenRegion, excess)
					},
				); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	check := func(_, right node.InspectableNode) (bool, string, hcl.Range) {
		return opt == isolateMultiLiners && right.Lines() > 1, separateMultiLine, right.Range()
	}

	_, err := r.requireLineBetween(runner, src, check, prev, cur)
	if err != nil {
		return err
	}

	check = func(left, right node.InspectableNode) (bool, string, hcl.Range) {
		return opt == isolateMultiLiners && left.Lines() > 1 && right.Lines() == 1, separateSingleLine, right.Range()
	}

	_, err = r.requireLineBetween(runner, src, check, prev, cur)
	if err != nil {
		return err
	}

	return nil
}

// checkTrailingSpace inspects the empty space after the last node inside the
// parent's range.
func (r *SpacingRule) checkTrailingSpace(
	runner *custom.Runner,
	src []byte,
	child node.Node,
	parent node.Node,
) error {
	lines := parent.Range().End.Line - child.Range().End.Line
	if lines <= 1 {
		return nil
	}

	_src := src[child.Range().End.Byte:parent.Range().End.Byte]

	_rng := hcl.Range{
		Filename: child.Range().Filename,
		Start:    child.Range().End,
		End:      parent.Range().End,
	}

	return r.checkLines(runner, src, _src, _rng, func(lines int, lastToken hcl.Range, sawComments bool) error {
		threshold := 0
		if sawComments {
			threshold = 1
		}

		if lines > threshold {
			excess := lines - threshold
			betweenRegion := hcl.Range{
				Filename: child.Range().Filename,
				Start:    child.Range().End,
				End:      parent.Range().End,
			}

			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf(
					"%d redundant blank %s in front",
					excess,
					lineOrLines(excess),
				),
				lastToken,
				func(fixer tflint.Fixer) error {
					return removeBlankLinesInRegion(fixer, src, betweenRegion, excess)
				},
			); err != nil {
				return err
			}
		}

		return nil
	})
}

// checkLines inspects the content within a range from the blank-line
// perspective (taking into consideration the comments).
func (r *SpacingRule) checkLines(
	runner *custom.Runner,
	fullSrc []byte,
	src []byte,
	rng hcl.Range,
	check func(lines int, lastToken hcl.Range, sawComments bool) error,
) error {
	tokens, err := hclsyntax.LexConfig(src, rng.Filename, rng.Start)
	if err != nil {
		return err
	}

	newLines := 0
	sawComments := false

	var (
		lastToken    hcl.Range
		lastTokenEnd hcl.Pos
	)

	lastTokenEnd = rng.Start

	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenNewline:
			newLines++
		case hclsyntax.TokenComment:
			lastToken = token.Range
			sawComments = true

			if newLines > 2 {
				excess := newLines - 2

				betweenRegion := hcl.Range{
					Filename: rng.Filename,
					Start:    lastTokenEnd,
					End:      token.Range.Start,
				}
				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf(
						"%d redundant blank %s in front",
						excess,
						lineOrLines(excess),
					),
					token.Range,
					func(fixer tflint.Fixer) error {
						return removeBlankLinesInRegion(fixer, fullSrc, betweenRegion, excess)
					},
				); err != nil {
					return err
				}
			}

			lastTokenEnd = token.Range.End
			newLines = 1 // Comment token implies `\n` at the end
		case hclsyntax.TokenEOF:
			// noop
		default:
			lastToken = token.Range
			lastTokenEnd = token.Range.End
		}
	}

	return check(newLines-1, lastToken, sawComments)
}

func (r *SpacingRule) requireLineBetween(
	runner *custom.Runner,
	src []byte,
	check func(left, right node.InspectableNode) (bool, string, hcl.Range),
	left, right node.InspectableNode,
) (bool, error) {
	lines := right.Range().Start.Line - left.Range().End.Line
	if emit, msg, issueRange := check(left, right); emit && lines < 2 {
		err := runner.EmitIssueWithFix(r, msg, issueRange,
			func(fixer tflint.Fixer) error {
				return insertBlankLineAfter(fixer, src, left.Range())
			},
		)

		return true, err
	}

	return false, nil
}

// checkModuleSpacing handles spacing for module blocks. When a module has
// key-attributes configured, source and key-attributes are grouped together
// (no blank line), with a blank line separating them from the rest. When no
// key-attributes are configured, a blank line is required after the source.
func (r *SpacingRule) checkModuleSpacing(
	runner *custom.Runner,
	block *hclsyntax.Block,
	src []byte,
	nodes []node.InspectableNode,
) ([]node.InspectableNode, error) {
	sourceStr := getSource(block)
	keys := keyAttrSet(runner.Modules, sourceStr)

	// Group source + for_each/count + key-attributes together.
	keys["source"] = true
	keys["for_each"] = true
	keys["count"] = true

	return r.checkKeyAttributeSpacing(runner, src, nodes, keys)
}

// checkKeyAttributeSpacing enforces no blank lines between consecutive
// key-attributes and a blank line separating them from the remaining attributes.
// When trim is true, processed key-attribute nodes are removed from the returned
// slice (safe when they form a contiguous prefix, e.g. modules). When false,
// nodes are returned as-is (for resources where key-attrs may follow count/for_each).
func (r *SpacingRule) checkKeyAttributeSpacing(
	runner *custom.Runner,
	src []byte,
	nodes []node.InspectableNode,
	keyAttrSet map[string]bool,
) ([]node.InspectableNode, error) {
	if disabled, exists := runner.Disabled["separate_key_attributes"]; exists && disabled {
		return nodes, nil
	}

	if len(keyAttrSet) == 0 || len(nodes) < 2 {
		return nodes, nil
	}

	// Find the first key-attribute.
	start := 0
	for start < len(nodes) {
		a := nodes[start].AsAttribute()
		if a != nil && keyAttrSet[a.Name] {
			break
		}

		start++
	}

	if start >= len(nodes) {
		return nodes, nil
	}

	// Walk over consecutive key-attributes, forbidding blank lines between them.
	// Only emit one issue per pass to avoid overlapping fix ranges.
	prev := nodes[start]
	i := start + 1
	emittedForbid := false

	for i < len(nodes) {
		cur := nodes[i]

		attribute := cur.AsAttribute()
		if attribute == nil || !keyAttrSet[attribute.Name] {
			break
		}

		ok, err := r.forbidLineBetween(runner, src, prev, cur, fmt.Sprintf(noSeparateKeyAttribute, attribute.Name))
		if err != nil {
			return nodes[i:], err
		}

		if ok {
			emittedForbid = true
		}

		prev = cur
		i++
	}

	// Require a blank line after the last key-attribute before the remaining
	// attributes. Skip if we already emitted a forbid fix to avoid
	// overlapping rewrite ranges.
	if !emittedForbid && i < len(nodes) {
		check := func(left, _ node.InspectableNode) (bool, string, hcl.Range) {
			return true, fmt.Sprintf(separateAttribute, left.Name()), left.Range()
		}

		if _, err := r.requireLineBetween(runner, src, check, nodes[i-1], nodes[i]); err != nil {
			return nodes, err
		}
	}

	return nodes[i:], nil
}

// forbidLineBetween emits an issue and fix when there is a blank line between
// two nodes that should be adjacent. Returns true if an issue was emitted.
func (r *SpacingRule) forbidLineBetween(
	runner *custom.Runner,
	src []byte,
	left, right node.InspectableNode,
	msg string,
) (bool, error) {
	lines := right.Range().Start.Line - left.Range().End.Line
	if lines <= 1 {
		return false, nil
	}

	excess := lines - 1
	betweenRegion := hcl.Range{
		Filename: left.Range().Filename,
		Start:    left.Range().End,
		End:      right.Range().Start,
	}
	err := runner.EmitIssueWithFix(r, msg, right.Range(),
		func(fixer tflint.Fixer) error {
			return removeBlankLinesInRegion(fixer, src, betweenRegion, excess)
		},
	)

	return true, err
}
