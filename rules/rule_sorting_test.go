package rules

import (
	"testing"
)

func TestSortingRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewSortingRule())
}
