package rules

import (
	"testing"
)

func TestLifecycleRule(t *testing.T) {
	t.Parallel()
	runTests(t, NewLifecycleRule())
}
