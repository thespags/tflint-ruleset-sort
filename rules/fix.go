package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/node"
)

// nodeLineRange returns a range covering the full lines of a node,
// including leading indentation and trailing newline.
func nodeLineRange(src []byte, nodeRange hcl.Range) hcl.Range {
	startByte := nodeRange.Start.Byte
	for startByte > 0 && src[startByte-1] != '\n' {
		startByte--
	}

	endByte := nodeRange.End.Byte
	for endByte < len(src) && src[endByte] != '\n' {
		endByte++
	}

	if endByte < len(src) {
		endByte++ // include the newline
	}

	return hcl.Range{
		Filename: nodeRange.Filename,
		Start:    hcl.Pos{Byte: startByte, Line: nodeRange.Start.Line, Column: 1},
		End:      hcl.Pos{Byte: endByte, Line: nodeRange.End.Line + 1, Column: 1},
	}
}

// moveNodeBefore moves a node (by its full lines) to before the target node's lines.
func moveNodeBefore(fixer tflint.Fixer, src []byte, toMove, target hcl.Range) error {
	moveRange := nodeLineRange(src, toMove)
	targetRange := nodeLineRange(src, target)
	moveText := string(src[moveRange.Start.Byte:moveRange.End.Byte])

	if err := fixer.Remove(moveRange); err != nil {
		return err
	}

	return fixer.InsertTextBefore(targetRange, moveText)
}

// moveNodeAfter moves a node (by its full lines) to after the target node's lines.
func moveNodeAfter(fixer tflint.Fixer, src []byte, toMove, target hcl.Range) error {
	moveRange := nodeLineRange(src, toMove)
	targetRange := nodeLineRange(src, target)
	moveText := string(src[moveRange.Start.Byte:moveRange.End.Byte])

	if err := fixer.Remove(moveRange); err != nil {
		return err
	}

	return fixer.InsertTextAfter(targetRange, moveText)
}

// reorderNodes replaces a consecutive run of nodes with the given new order
// in a single ReplaceText call. `oldOrder` is the nodes as they appear in the
// source; `newOrder` is the desired sequence. Spacing between positions is
// preserved (i.e., the gap between position 0 and 1 stays between whatever
// nodes now occupy positions 0 and 1).
func reorderNodes(fixer tflint.Fixer, src []byte, oldOrder, newOrder []node.InspectableNode) error {
	// Build replacement: new-order texts with original-position spacings.
	var (
		replacement string
		sb          strings.Builder
	)

	for i, n := range newOrder {
		sb.Write(src[n.Range().Start.Byte:n.Range().End.Byte])

		if i < len(newOrder)-1 {
			// Preserve the spacing that was between position i and i+1.
			sb.Write(src[oldOrder[i].Range().End.Byte:oldOrder[i+1].Range().Start.Byte])
		}
	}

	replacement += sb.String()

	span := hcl.RangeBetween(oldOrder[0].Range(), oldOrder[len(oldOrder)-1].Range())

	return fixer.ReplaceText(span, replacement)
}

// insertBlankLineAfter inserts a blank line after the given range's line.
func insertBlankLineAfter(fixer tflint.Fixer, src []byte, r hcl.Range) error {
	lr := nodeLineRange(src, r)

	return fixer.InsertTextAfter(lr, "\n")
}

// removeBlankLinesInRegion removes n consecutive blank lines from within a
// byte region. It looks for consecutive \n\n and removes one \n per blank line.
func removeBlankLinesInRegion(fixer tflint.Fixer, src []byte, region hcl.Range, n int) error {
	if region.Start.Byte >= region.End.Byte || n <= 0 {
		return nil
	}

	between := string(src[region.Start.Byte:region.End.Byte])

	newBetween := between

	removed := 0
	for i := 0; i < len(newBetween)-1 && removed < n; i++ {
		if newBetween[i] == '\n' && i+1 < len(newBetween) && newBetween[i+1] == '\n' {
			newBetween = newBetween[:i+1] + newBetween[i+2:]
			removed++
			i-- // re-check same position
		}
	}

	if removed == 0 {
		return nil
	}

	return fixer.ReplaceText(region, newBetween)
}
