package custom

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thespags/tflint-ruleset-sort/config"
)

// TestLookup_OneWayFallback exercises the data → resource fallback,
// and verifies the resource → data direction does NOT fall back.
func TestLookup_OneWayFallback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		resources   []*config.Resource
		dataBlocks  []*config.Resource
		blockType   string
		kind        string
		wantNil     bool
		wantKeyAttr []string
	}{
		{
			name:        "data block hits its own config",
			dataBlocks:  []*config.Resource{{Kind: "thing", Keys: []string{"d"}}},
			blockType:   "data",
			kind:        "thing",
			wantKeyAttr: []string{"d"},
		},
		{
			name:        "data block falls back to resource config",
			resources:   []*config.Resource{{Kind: "thing", Keys: []string{"r"}}},
			blockType:   "data",
			kind:        "thing",
			wantKeyAttr: []string{"r"},
		},
		{
			name:        "data block prefers data over resource when both defined",
			resources:   []*config.Resource{{Kind: "thing", Keys: []string{"r"}}},
			dataBlocks:  []*config.Resource{{Kind: "thing", Keys: []string{"d"}}},
			blockType:   "data",
			kind:        "thing",
			wantKeyAttr: []string{"d"},
		},
		{
			name:        "resource block hits its own config",
			resources:   []*config.Resource{{Kind: "thing", Keys: []string{"r"}}},
			blockType:   "resource",
			kind:        "thing",
			wantKeyAttr: []string{"r"},
		},
		{
			name:       "resource block does NOT fall back to data config",
			dataBlocks: []*config.Resource{{Kind: "thing", Keys: []string{"d"}}},
			blockType:  "resource",
			kind:       "thing",
			wantNil:    true,
		},
		{
			name:      "neither configured returns nil",
			blockType: "data",
			kind:      "missing",
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{Resources: tt.resources, Data: tt.dataBlocks}

			runner, err := NewRunner(nil, cfg)
			require.NoError(t, err)

			block := &hclsyntax.Block{Type: tt.blockType, Labels: []string{tt.kind, "instance_name"}}

			got := runner.Lookup(block)
			if tt.wantNil {
				assert.Nil(t, got)

				return
			}

			require.NotNil(t, got)
			assert.Equal(t, tt.wantKeyAttr, got.KeyAttributes)
		})
	}
}
