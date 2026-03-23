package custom

import (
	"testing"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigResource(t *testing.T) {
	t.Parallel()

	res, err := parseConfigResource(&config.Resource{
		Kind: "kubernetes_manifest",
		Keys: []string{"manifest.metadata.namespace", "manifest.metadata.name"},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"manifest", "metadata"}, res.KeyBlocks)
	assert.Equal(t, []string{"namespace", "name"}, res.KeyAttributes)
}
