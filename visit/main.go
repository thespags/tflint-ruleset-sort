package visit

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Files visits all files in a runner.
func Files(runner tflint.Runner, visit func(*hclsyntax.Body, []byte) error) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		// `.tf.json` files have a JSON body; this ruleset's sorting/spacing
		// rules only apply to HCL source layout, so skip non-HCL files.
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		if err := visit(body, file.Bytes); err != nil {
			return err
		}
	}

	return nil
}

// Blocks visits all blocks in a file.
func Blocks(runner tflint.Runner, visit func(*hclsyntax.Block, []byte) error) error {
	return Files(runner, func(body *hclsyntax.Body, bytes []byte) error {
		for _, block := range body.Blocks {
			if err := visit(block, bytes); err != nil {
				return err
			}
		}

		return nil
	})
}
