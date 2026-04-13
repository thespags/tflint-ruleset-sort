package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/thespags/tflint-ruleset-sort/custom"
	"github.com/thespags/tflint-ruleset-sort/node"
	"github.com/zclconf/go-cty/cty"
)

func toNames(nodes []node.InspectableNode) string {
	sb := strings.Builder{}
	for _, n := range nodes {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(n.Name())
	}

	return sb.String()
}

func nameOrdered(a, b node.InspectableNode) int {
	return strings.Compare(a.Name(), b.Name())
}

func singleLinesFirst(a, b node.InspectableNode) int {
	if a.Lines() > 1 && b.Lines() == 1 {
		return 1
	}

	return -1
}

func attributesFirst(a, b node.InspectableNode) int {
	if a.IsBlock() && b.IsAttribute() {
		return 1
	}

	return -1
}

func sortedRange(oldOrder, newOrder []node.InspectableNode) hcl.Range {
	start := oldOrder[0]
	for i := range oldOrder {
		if oldOrder[i] != newOrder[i] {
			start = oldOrder[i]

			break
		}
	}

	end := oldOrder[len(oldOrder)-1]
	for i := len(oldOrder) - 1; i >= 0; i-- {
		if oldOrder[i] != newOrder[i] {
			end = oldOrder[i]

			break
		}
	}

	return hcl.RangeBetween(start.Range(), end.Range())
}

func hasCommentsInBetween(src []byte, left node.Node, right node.Node) (bool, error) {
	rng := hcl.Range{
		Filename: right.Range().Filename,
		Start: hcl.Pos{
			Line:   left.Range().End.Line,
			Column: left.Range().End.Column,
			Byte:   left.Range().End.Byte,
		},
		End: hcl.Pos{
			Line:   right.Range().Start.Line,
			Column: right.Range().Start.Column,
			Byte:   right.Range().Start.Byte,
		},
	}

	tokens, err := hclsyntax.LexConfig(
		src[rng.Start.Byte:rng.End.Byte],
		rng.Filename,
		rng.Start,
	)
	if err != nil {
		return false, err
	}

	for _, token := range tokens {
		if token.Type == hclsyntax.TokenComment {
			return true, nil
		}
	}

	return false, nil
}

func lineOrLines(count int) string {
	if count == 1 {
		return "line"
	}

	return "lines"
}

func getSource(block *hclsyntax.Block) string {
	source, exists := block.Body.Attributes["source"]
	if !exists {
		return ""
	}

	val, diags := source.Expr.Value(nil)
	if diags.HasErrors() || val.Type() != cty.String {
		return ""
	}

	return val.AsString()
}

// resourceKeyAttrSet returns the set of key-attribute names for a resource/data block,
// or nil if none are configured or if key-attributes are nested in key-blocks.
func keyAttrSet(resources map[string]*custom.Resource, kind string) map[string]bool {
	resource, known := resources[kind]
	if !known || len(resource.KeyAttributes) == 0 {
		return make(map[string]bool)
	}

	set := make(map[string]bool, len(resource.KeyAttributes)+1)
	for _, ka := range resource.KeyAttributes {
		set[ka] = true
	}

	return set
}
