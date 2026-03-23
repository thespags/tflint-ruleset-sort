package rules

import (
	"testing"
)

func TestKeyAttributesRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewKeyAttributesRule())
}
