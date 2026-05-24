package visit

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestFiles_SkipsJSON(t *testing.T) {
	t.Parallel()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `resource "aws_s3_bucket" "b" {
  bucket = "example"
}
`,
		"generated.tf.json": `{
  "resource": {
    "aws_s3_bucket": {
      "b": { "bucket": "example" }
    }
  }
}
`,
	})

	visited := map[string]bool{}
	err := Files(runner, func(body *hclsyntax.Body, _ []byte) error {
		for _, block := range body.Blocks {
			visited[block.Type] = true
		}

		return nil
	})

	require.NoError(t, err, "Files should skip .tf.json without erroring")
	require.True(t, visited["resource"], "HCL file should still have been visited")
}
