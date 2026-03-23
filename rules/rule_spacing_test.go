package rules

import (
	"testing"
)

func TestSpacingRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewSpacingRule())
}
