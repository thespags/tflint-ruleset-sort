package node

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type block struct {
	block *hclsyntax.Block
}

// WrapBlock wraps an HCL block as an InspectableNode.
func WrapBlock(hclBlock *hclsyntax.Block) InspectableNode {
	return &block{
		block: hclBlock,
	}
}

// Range returns the source range of the block.
func (b block) Range() hcl.Range {
	return b.block.Range()
}

// Kind returns Block.
func (block) Kind() Kind {
	return Block
}

// AsAttribute returns nil (blocks are not attributes).
func (block) AsAttribute() *hclsyntax.Attribute {
	return nil
}

// AsBlock returns the underlying HCL block.
func (b block) AsBlock() *hclsyntax.Block {
	return b.block
}

// IsAttribute returns false.
func (block) IsAttribute() bool {
	return false
}

// IsBlock returns true.
func (block) IsBlock() bool {
	return true
}

// Name returns the block type and labels joined by spaces.
// For dynamic blocks, the type prefix is omitted.
func (b block) Name() string {
	names := make([]string, 0, len(b.block.Labels)+1)
	if b.block.Type != "dynamic" {
		names = append(names, b.block.Type)
	}

	names = append(names, b.block.Labels...)

	return strings.Join(names, " ")
}

// Type returns "block".
func (block) Type() string {
	return "block"
}

// Lines returns the number of lines the block spans.
func (b block) Lines() int {
	r := b.Range()

	return r.End.Line - r.Start.Line + 1
}
