package rules

import (
	"testing"
)

func TestUnknownResourceRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewUnknownResourceRule())
}
