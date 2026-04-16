package node

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type objectConsItem struct {
	item *hclsyntax.ObjectConsItem
}

// WrapObjectConsItem wraps an HCL object key-value pair as an InspectableNode.
func WrapObjectConsItem(i *hclsyntax.ObjectConsItem) InspectableNode {
	return &objectConsItem{
		item: i,
	}
}

// Range returns the source range spanning from the key to the value expression.
func (i objectConsItem) Range() hcl.Range {
	return hcl.RangeBetween(
		i.item.KeyExpr.Range(),
		i.item.ValueExpr.Range(),
	)
}

// Kind returns Composite.
func (objectConsItem) Kind() Kind {
	return Composite
}

// AsAttribute returns nil (object items are not attributes).
func (objectConsItem) AsAttribute() *hclsyntax.Attribute {
	return nil
}

// AsBlock returns nil (object items are not blocks).
func (objectConsItem) AsBlock() *hclsyntax.Block {
	return nil
}

// IsAttribute returns false.
func (objectConsItem) IsAttribute() bool {
	return false
}

// IsBlock returns false.
func (objectConsItem) IsBlock() bool {
	return false
}

// Type returns "key".
func (objectConsItem) Type() string {
	return "key"
}

// Name returns the key name extracted from the key expression.
func (i objectConsItem) Name() string {
	var names []string

	if x, ok := i.item.KeyExpr.(*hclsyntax.ObjectConsKeyExpr); ok {
		switch expr := x.Wrapped.(type) {
		case *hclsyntax.ScopeTraversalExpr:
			for _, traversal := range expr.Traversal {
				names = make([]string, 0, len(expr.Traversal))

				switch typedTravesrsal := traversal.(type) {
				// TODO: Figure out other possible types
				case hcl.TraverseRoot:
					names = append(names, typedTravesrsal.Name)
				case hcl.TraverseAttr:
					names = append(names, typedTravesrsal.Name)
				case *hcl.TraverseRoot:
					names = append(names, typedTravesrsal.Name)
				case *hcl.TraverseAttr:
					names = append(names, typedTravesrsal.Name)
				}
			}
		case *hclsyntax.TemplateExpr:
			names = make([]string, 0, len(expr.Parts))
			for _, part := range expr.Parts {
				if typedPart, ok := part.(*hclsyntax.LiteralValueExpr); ok {
					names = append(names, typedPart.Val.AsString())
				}
			}
		}
	}

	pos := 0

	for _, n := range names {
		if len(n) > 0 {
			names[pos] = n
			pos++
		}
	}

	names = names[:pos]

	return strings.Join(names, ".")
}

// Expr returns the value expression of the object item.
func (i objectConsItem) Expr() hclsyntax.Expression {
	return i.item.ValueExpr
}

// Lines returns the number of lines the object item spans.
func (i objectConsItem) Lines() int {
	r := i.Range()

	return r.End.Line - r.Start.Line + 1
}
