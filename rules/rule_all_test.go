package rules

import (
	"testing"
)

func TestAllRules(t *testing.T) {
	t.Parallel()
	runTests(t, All()...)
}
