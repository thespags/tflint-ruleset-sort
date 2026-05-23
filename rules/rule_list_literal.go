package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/custom"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// ListLiteralRule sorts the elements of list literals (HCL `TupleConsExpr`)
// everywhere they appear — top-level attribute values, function-call args,
// nested object values, `for` expressions, etc. Per-resource opt-out via
// `skip_sort_literals` for attributes whose list order is semantically
// meaningful (`command`, `args`, URL `path` matchers, ...).
type ListLiteralRule struct {
	tflint.DefaultRule
}

// NewListLiteralRule creates a new ListLiteralRule.
func NewListLiteralRule() *ListLiteralRule {
	return &ListLiteralRule{}
}

// Name returns the name of the rule.
func (*ListLiteralRule) Name() string {
	return project.RuleName("list_literal")
}

// Enabled returns whether the rule is enabled by default.
func (*ListLiteralRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*ListLiteralRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *ListLiteralRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check walks the AST and sorts list-literal elements (subject to opt-out).
func (r *ListLiteralRule) Check(rr tflint.Runner) error {
	runner, ok := rr.(*custom.Runner)
	if !ok {
		return nil
	}

	return visit.Files(runner, func(body *hclsyntax.Body, src []byte) error {
		return r.walkBody(runner, src, body, nil)
	})
}

// walkBody walks a body's attributes and nested blocks. `skipLiterals` (may
// be nil) is the active set of attribute names whose value subtrees must not
// have their list-literal elements sorted.
func (r *ListLiteralRule) walkBody(
	runner *custom.Runner,
	src []byte,
	body *hclsyntax.Body,
	skipLiterals map[string]bool,
) error {
	for _, attr := range body.Attributes {
		if err := r.walkExpr(runner, src, attr.Expr, skipLiterals[attr.Name]); err != nil {
			return err
		}
	}

	for _, block := range body.Blocks {
		var inner map[string]bool
		if cfg := runner.Lookup(block); cfg != nil && len(cfg.SkipSortLiterals) > 0 {
			inner = make(map[string]bool, len(cfg.SkipSortLiterals))
			for _, name := range cfg.SkipSortLiterals {
				inner[name] = true
			}
		}

		if err := r.walkBody(runner, src, block.Body, inner); err != nil {
			return err
		}
	}

	return nil
}

// walkExpr recurses through an expression tree. When `skip` is true the
// current tuple is not sorted, and the flag propagates to its descendants —
// nested lists inside a skipped attribute are also left alone.
func (r *ListLiteralRule) walkExpr(
	runner *custom.Runner,
	src []byte,
	expression hclsyntax.Expression,
	skip bool,
) error {
	switch expr := expression.(type) {
	case *hclsyntax.TupleConsExpr:
		if !skip {
			if err := r.sortTupleElements(runner, src, expr); err != nil {
				return err
			}
		}

		for _, sub := range expr.Exprs {
			if err := r.walkExpr(runner, src, sub, skip); err != nil {
				return err
			}
		}

	case *hclsyntax.ObjectConsExpr:
		for _, item := range expr.Items {
			if err := r.walkExpr(runner, src, item.ValueExpr, skip); err != nil {
				return err
			}
		}

	case *hclsyntax.FunctionCallExpr:
		for _, arg := range expr.Args {
			if err := r.walkExpr(runner, src, arg, skip); err != nil {
				return err
			}
		}

	case *hclsyntax.ForExpr:
		return r.walkExpr(runner, src, expr.ValExpr, skip)

	case *hclsyntax.ParenthesesExpr:
		return r.walkExpr(runner, src, expr.Expression, skip)
	}

	return nil
}

// sortTupleElements emits an issue and fix when the tuple's elements are not in
// natural-numeric-aware order. Lists containing expressions that aren't
// lexically orderable (function calls, complex interpolations) are silently
// skipped. Trailing inline comments on each element travel with the element,
// as do contiguous `//` / `#` comment lines directly above each element.
func (r *ListLiteralRule) sortTupleElements(
	runner *custom.Runner,
	src []byte,
	tuple *hclsyntax.TupleConsExpr,
) error {
	if len(tuple.Exprs) < 2 {
		return nil
	}

	names := make([]string, len(tuple.Exprs))

	for i, e := range tuple.Exprs {
		name, ok := elementName(src, e)
		if !ok {
			return nil
		}

		names[i] = name
	}

	order := make([]int, len(tuple.Exprs))
	for i := range order {
		order[i] = i
	}

	slices.SortStableFunc(order, func(a, b int) int {
		return natCompare(names[a], names[b])
	})

	if isIdentityOrder(order) {
		return nil
	}

	rebuilt := rebuildTuple(src, tuple, order)

	firstWrong := 0

	for i, o := range order {
		if i != o {
			firstWrong = i

			break
		}
	}

	wrongRange := tuple.Exprs[firstWrong].Range()

	return runner.EmitIssueWithFix(
		r,
		fmt.Sprintf(
			"list-literal is not sorted: expected `%s` at position %d, got `%s`",
			names[order[firstWrong]],
			firstWrong,
			names[firstWrong],
		),
		wrongRange,
		func(fixer tflint.Fixer) error {
			return fixer.ReplaceText(tuple.Range(), rebuilt)
		},
	)
}

func isIdentityOrder(order []int) bool {
	for i, o := range order {
		if i != o {
			return false
		}
	}

	return true
}

// elementName returns the comparable name of a tuple element, and false if
// the expression cannot be lexically ordered (function calls, computed
// interpolations, etc.).
func elementName(src []byte, expression hclsyntax.Expression) (string, bool) {
	switch expr := expression.(type) {
	case *hclsyntax.LiteralValueExpr:
		r := expr.Range()

		return string(src[r.Start.Byte:r.End.Byte]), true
	case *hclsyntax.TemplateExpr:
		if !expr.IsStringLiteral() {
			return "", false
		}

		val, diags := expr.Value(nil)
		if diags.HasErrors() {
			return "", false
		}

		return val.AsString(), true
	case *hclsyntax.TemplateWrapExpr:
		return elementName(src, expr.Wrapped)
	case *hclsyntax.ScopeTraversalExpr:
		parts := make([]string, 0, len(expr.Traversal))

		for _, t := range expr.Traversal {
			switch tt := t.(type) {
			case hcl.TraverseRoot:
				parts = append(parts, tt.Name)
			case hcl.TraverseAttr:
				parts = append(parts, tt.Name)
			default:
				return "", false
			}
		}

		return strings.Join(parts, "."), true
	default:
		return "", false
	}
}

func rebuildTuple(src []byte, tuple *hclsyntax.TupleConsExpr, order []int) string {
	multiLine := tuple.Exprs[0].Range().Start.Line != tuple.OpenRange.Start.Line

	if !multiLine {
		parts := make([]string, len(order))
		for newIdx, oldIdx := range order {
			r := tuple.Exprs[oldIdx].Range()
			parts[newIdx] = string(src[r.Start.Byte:r.End.Byte])
		}

		return "[" + strings.Join(parts, ", ") + "]"
	}

	tupleStart := tuple.Range().Start.Byte + 1
	tupleEnd := tuple.Range().End.Byte - 1

	parts := make([]string, len(tuple.Exprs))
	lineStart := nextNewLine(src, tupleStart, tupleEnd) + 1
	prefix := string(src[tupleStart:lineStart])

	for srcIdx, expr := range tuple.Exprs {
		rng := expr.Range()

		lineEnd := nextNewLine(src, rng.End.Byte, tupleEnd)
		parts[srcIdx] = string(src[lineStart:lineEnd])
		lineStart = lineEnd + 1
	}

	sorted := make([]string, len(parts))
	for newIdx, oldIdx := range order {
		sorted[newIdx] = parts[oldIdx]
	}

	suffix := string(src[lineStart:tupleEnd])

	return "[" + prefix + strings.Join(sorted, "\n") + "\n" + suffix + "]"
}

func nextNewLine(src []byte, start, end int) int {
	i := start
	for i < end && src[i] != '\n' {
		i++
	}

	return i
}
