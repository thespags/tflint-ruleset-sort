package custom

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/config"
	"github.com/zclconf/go-cty/cty"
)

// Runner is a wrapper of RPC client with custom configuration.
type Runner struct {
	tflint.Runner

	// Disabled tells whether a rule is disabled or not.
	Disabled map[string]bool

	// Resources stores the configuration of terraform `resource` blocks.
	Resources map[string]*Resource

	// Data stores the configuration of terraform `data` blocks.
	Data map[string]*Resource

	// Modules stores the configuration of terraform modules (keyed by source).
	Modules map[string]*Resource
}

// Resource is the configuration of the `resource` and `data` blocks that linter
// uses to apply its rules.
type Resource struct {
	// KeyBlocks is the (sequence of nested) block type(s) that contain
	// key-attributes (for example, `metadata` in kubernetes resources).
	KeyBlocks []string

	// KeyAttributes is the prioritized list of attributes that uniquely
	// identify the `resource` or `data` block.
	KeyAttributes []string
}

// NewRunner returns a new runner.
func NewRunner(runner tflint.Runner, customConfig *config.Config) (*Runner, error) {
	resources := map[string]*Resource{}

	for _, resource := range customConfig.Resources {
		res, err := parseConfigResource(resource)
		if err != nil {
			return nil, err
		}

		resources[resource.Kind] = res
	}

	dataBlocks := map[string]*Resource{}

	for _, dataBlock := range customConfig.Data {
		res, err := parseConfigResource(dataBlock)
		if err != nil {
			return nil, err
		}

		dataBlocks[dataBlock.Kind] = res
	}

	modules := map[string]*Resource{}

	for _, module := range customConfig.Modules {
		res, err := parseConfigResource(module)
		if err != nil {
			return nil, err
		}

		modules[module.Kind] = res
	}

	return &Runner{
		Runner:    runner,
		Disabled:  make(map[string]bool),
		Resources: resources,
		Data:      dataBlocks,
		Modules:   modules,
	}, nil
}

// Lookup returns the configured Resource for the given block, or nil if not
// configured. It dispatches by `block.Type`:
//
//   - `resource` blocks → keyed by `block.Labels[0]` against Resources.
//   - `data` blocks     → keyed by `block.Labels[0]` against Data, falling
//     back to Resources if not present (a data source is conceptually a
//     read-only view of the resource of the same name). The reverse direction
//     does not fall back — open to change but left out for less magic.
//   - `module` blocks   → keyed by the block's `source` attribute value
//     against Modules.
//
// Any other block type returns nil.
func (r *Runner) Lookup(block *hclsyntax.Block) *Resource {
	switch block.Type {
	case "data":
		if res, ok := r.Data[block.Labels[0]]; ok {
			return res
		}
		return r.Resources[block.Labels[0]]
	case "resource":
		return r.Resources[block.Labels[0]]
	case "module":
		return r.Modules[GetSource(block)]
	}

	return nil
}

// GetSource returns the literal value of a module block's `source`
// attribute, or an empty string if it is missing or non-string.
func GetSource(block *hclsyntax.Block) string {
	source, exists := block.Body.Attributes["source"]
	if !exists {
		return ""
	}

	val, diags := source.Expr.Value(nil)
	if diags.HasErrors() || val.Type() != cty.String {
		return ""
	}

	return val.AsString()
}
