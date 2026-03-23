package node

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type expression struct {
	expression hclsyntax.Expression
}

// WrapExpression wraps an HCL expression as an InspectableNode.
func WrapExpression(a hclsyntax.Expression) InspectableNode {
	return &expression{
		expression: a,
	}
}

// Range returns the source range of the expression.
func (x expression) Range() hcl.Range {
	return x.expression.Range()
}

// Kind returns Expression.
func (expression) Kind() Kind {
	return Expression
}

// AsAttribute returns nil (expressions are not attributes).
func (expression) AsAttribute() *hclsyntax.Attribute {
	return nil
}

// AsBlock returns nil (expressions are not blocks).
func (expression) AsBlock() *hclsyntax.Block {
	return nil
}

// IsAttribute returns true (expressions are treated as attribute-like for sorting).
func (expression) IsAttribute() bool {
	return true
}

// IsBlock returns false.
func (expression) IsBlock() bool {
	return false
}

// Name panics because expressions have no meaningful name.
func (expression) Name() string {
	panic("expression has no name")
}

// Type returns "expression".
func (expression) Type() string {
	return "expression"
}

// Lines returns the number of lines the expression spans.
func (x expression) Lines() int {
	r := x.Range()

	return r.End.Line - r.Start.Line + 1
}
