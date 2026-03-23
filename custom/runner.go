package custom

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Runner is a wrapper of RPC client with custom configuration.
type Runner struct {
	tflint.Runner

	// Disabled tells whether a rule is disabled or not.
	Disabled map[string]bool

	// Resources stores the configuration of terraform resources.
	Resources map[string]*Resource

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
		Modules:   modules,
	}, nil
}
