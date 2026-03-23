package rules

import (
	"testing"
)

func TestCountRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewCountRule())
}
