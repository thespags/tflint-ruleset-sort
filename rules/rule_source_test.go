package rules

import (
	"testing"
)

func TestSourceRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewSourceRule())
}
