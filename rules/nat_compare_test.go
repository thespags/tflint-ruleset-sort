package rules

import (
	"testing"
)

func TestNatCompare(t *testing.T) {
	t.Parallel()

	cases := []struct {
		a, b string
		want int // sign only: -1, 0, +1
	}{
		{"a", "b", -1},
		{"b", "a", +1},
		{"a", "a", 0},
		{"v1", "v2", -1},
		{"v2", "v10", -1},
		{"v10", "v2", +1},
		{"v10", "v10", 0},
		{"3", "100", -1},
		{"100", "3", +1},
		{"foo", "100", +1}, // digits ('1' = 0x31) sort before letters ('f' = 0x66)
		{"100", "foo", -1},
		{"adas36", "afailla", -1},
		{"item1", "item01", -1}, // shorter original wins on numeric tie
		{"", "a", -1},
		{"a", "", +1},
	}

	for _, c := range cases {
		got := natCompare(c.a, c.b)
		gotSign := sign(got)

		if gotSign != c.want {
			t.Errorf("natCompare(%q, %q) = %d (sign %d); want sign %d",
				c.a, c.b, got, gotSign, c.want)
		}
	}
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return +1
	default:
		return 0
	}
}
