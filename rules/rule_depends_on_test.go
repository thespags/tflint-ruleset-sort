package rules

import (
	"testing"
)

func TestDependsOnRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewDependsOnRule())
}
