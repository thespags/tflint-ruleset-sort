package rules

import (
	"testing"
)

func TestForEachRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewForEachRule())
}
