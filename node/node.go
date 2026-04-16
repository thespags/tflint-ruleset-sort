package node

import (
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Node is a drop-in replacement for HCL's nodes.
type Node interface {
	Range() hcl.Range
}

// InspectableNode provides additional information about the node.
type InspectableNode interface {
	Node

	// AsAttribute returns the underlying HCL attribute, or nil if the node is not an attribute.
	AsAttribute() *hclsyntax.Attribute

	// AsBlock returns the underlying HCL block, or nil if the node is not a block.
	AsBlock() *hclsyntax.Block

	// IsAttribute reports whether the node is an attribute.
	IsAttribute() bool

	// IsBlock reports whether the node is a block.
	IsBlock() bool

	// Lines returns the number of lines the node spans (1 for single-line).
	Lines() int

	// Kind returns the kind of the node (Attribute, Block, Expression, or Composite).
	Kind() Kind

	// Name returns a human-readable name for the node.
	Name() string

	// Type returns a string label for the node kind (e.g. "attribute", "block").
	Type() string

	// Expr returns the expression to recurse into, or nil.
	Expr() hclsyntax.Expression
}

// Kind indicates the kind of a node (e.g. attribute, or block, or expression,
// and so on).
type Kind int

const (
	// Attribute is an HCL attribute (key = value).
	Attribute Kind = iota

	// Block is an HCL block (e.g. resource, lifecycle).
	Block

	// Expression is a standalone HCL expression.
	Expression

	// Composite node combines multiple nodes together.
	// For example, key + value expressions.
	Composite
)

// FirstNodeFrom returns the first node (by source position) from the body,
// or nil if the body is empty.
func FirstNodeFrom(body *hclsyntax.Body) InspectableNode {
	if len(body.Attributes) == 0 && len(body.Blocks) == 0 {
		return nil
	}

	var first *hclsyntax.Attribute
	for _, current := range body.Attributes {
		if first == nil {
			first = current

			continue
		}

		if first.SrcRange.Start.Byte < current.SrcRange.Start.Byte {
			continue
		}

		first = current
	}

	var hclBlock *hclsyntax.Block
	if len(body.Blocks) > 0 {
		hclBlock = body.Blocks[0]
	}

	for i := 1; i < len(body.Blocks); i++ {
		if hclBlock == nil {
			hclBlock = body.Blocks[i]

			continue
		}

		if hclBlock.TypeRange.Start.Byte < body.Blocks[i].TypeRange.Start.Byte {
			continue
		}

		hclBlock = body.Blocks[i]
	}

	if first == nil {
		return WrapBlock(hclBlock)
	}

	if hclBlock == nil {
		return WrapAttribute(first)
	}

	if first.SrcRange.Start.Byte < hclBlock.TypeRange.Start.Byte {
		return WrapAttribute(first)
	}

	return WrapBlock(hclBlock)
}

// LastNodeFrom returns the last node (by source position) from the body,
// or nil if the body is empty.
func LastNodeFrom(body *hclsyntax.Body) InspectableNode {
	if len(body.Attributes) == 0 && len(body.Blocks) == 0 {
		return nil
	}

	var last *hclsyntax.Attribute
	for _, current := range body.Attributes {
		if last == nil {
			last = current

			continue
		}

		if last.SrcRange.Start.Byte > current.SrcRange.Start.Byte {
			continue
		}

		last = current
	}

	var hclBlock *hclsyntax.Block
	if len(body.Blocks) > 0 {
		hclBlock = body.Blocks[0]
	}

	for i := 1; i < len(body.Blocks); i++ {
		if hclBlock == nil {
			hclBlock = body.Blocks[i]

			continue
		}

		if hclBlock.TypeRange.Start.Byte > body.Blocks[i].TypeRange.Start.Byte {
			continue
		}

		hclBlock = body.Blocks[i]
	}

	if last == nil {
		return WrapBlock(hclBlock)
	}

	if hclBlock == nil {
		return WrapAttribute(last)
	}

	if last.SrcRange.Start.Byte > hclBlock.TypeRange.Start.Byte {
		return WrapAttribute(last)
	}

	return WrapBlock(hclBlock)
}

// OrderedInspectableNodesFrom returns all attributes and blocks from the body
// as InspectableNodes, sorted by their source position.
func OrderedInspectableNodesFrom(body *hclsyntax.Body) []InspectableNode {
	nodes := make([]InspectableNode, 0, len(body.Blocks)+len(body.Attributes))

	for _, hclAttribute := range body.Attributes {
		nodes = append(nodes, WrapAttribute(hclAttribute))
	}

	for _, hclBlock := range body.Blocks {
		nodes = append(nodes, WrapBlock(hclBlock))
	}

	slices.SortStableFunc(nodes, func(left, right InspectableNode) int {
		return left.Range().Start.Byte - right.Range().Start.Byte
	})

	return nodes
}
