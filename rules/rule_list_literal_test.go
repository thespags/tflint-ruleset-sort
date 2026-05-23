package rules

import (
	"testing"
)

func TestListLiteralRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewListLiteralRule())
}
