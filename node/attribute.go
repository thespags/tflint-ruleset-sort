package node

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type attribute struct {
	attribute *hclsyntax.Attribute
}

// WrapAttribute wraps an HCL attribute as an InspectableNode.
func WrapAttribute(a *hclsyntax.Attribute) InspectableNode {
	return &attribute{
		attribute: a,
	}
}

// Range returns the source range of the attribute.
func (a attribute) Range() hcl.Range {
	return a.attribute.Range()
}

// Kind returns Attribute.
func (attribute) Kind() Kind {
	return Attribute
}

// AsAttribute returns the underlying HCL attribute.
func (a attribute) AsAttribute() *hclsyntax.Attribute {
	return a.attribute
}

// AsBlock returns nil (attributes are not blocks).
func (attribute) AsBlock() *hclsyntax.Block {
	return nil
}

// IsAttribute returns true.
func (attribute) IsAttribute() bool {
	return true
}

// IsBlock returns false.
func (attribute) IsBlock() bool {
	return false
}

// Name returns the attribute name.
func (a attribute) Name() string {
	return a.attribute.Name
}

// Type returns "attribute".
func (attribute) Type() string {
	return "attribute"
}

// Lines returns the number of lines the attribute spans.
func (a attribute) Lines() int {
	r := a.Range()

	return r.End.Line - r.Start.Line + 1
}
