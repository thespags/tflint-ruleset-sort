package rules

import "strings"

// natCompare compares two strings with numeric-aware ordering. Numeric value
// compares contiguous digit runs (so "v2" < "v10" and "3" < "100").
// Outside digit runs, comparison is byte-wise — non-ASCII text is ordered by
// its UTF-8 byte sequence, not by locale (an `ä` in Swedish input will not
// sort to its alphabet-correct position). Leading-zero runs of equal numeric
// value sort the shorter original first ("1" < "01" < "001").
func natCompare(left, right string) int {
	leftIndex, rightIndex := 0, 0

	for leftIndex < len(left) && rightIndex < len(right) {
		leftDigit := isDigit(left[leftIndex])
		rightDigit := isDigit(right[rightIndex])

		if leftDigit && rightDigit {
			leftEnd := leftIndex
			for leftEnd < len(left) && isDigit(left[leftEnd]) {
				leftEnd++
			}

			rightEnd := rightIndex
			for rightEnd < len(right) && isDigit(right[rightEnd]) {
				rightEnd++
			}

			leftStart := leftIndex
			for leftStart < leftEnd-1 && left[leftStart] == '0' {
				leftStart++
			}

			rightStart := rightIndex
			for rightStart < rightEnd-1 && right[rightStart] == '0' {
				rightStart++
			}

			if leftLen, rightLen := leftEnd-leftStart, rightEnd-rightStart; leftLen != rightLen {
				return leftLen - rightLen
			}

			if cmp := strings.Compare(left[leftStart:leftEnd], right[rightStart:rightEnd]); cmp != 0 {
				return cmp
			}

			if leftOrigLen, rightOrigLen := leftEnd-leftIndex, rightEnd-rightIndex; leftOrigLen != rightOrigLen {
				return leftOrigLen - rightOrigLen
			}

			leftIndex = leftEnd
			rightIndex = rightEnd

			continue
		}

		if left[leftIndex] != right[rightIndex] {
			return int(left[leftIndex]) - int(right[rightIndex])
		}

		leftIndex++
		rightIndex++
	}

	return len(left) - len(right)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
