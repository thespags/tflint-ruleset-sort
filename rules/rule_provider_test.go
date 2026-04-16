package rules

import (
	"testing"
)

func TestProviderRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewProviderRule())
}
