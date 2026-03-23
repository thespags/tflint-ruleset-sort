package custom

import (
	config2 "github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// RuleSet is the custom ruleset.
type RuleSet struct {
	tflint.BuiltinRuleSet

	config *config2.Config
}

// ConfigSchema returns the ruleset plugin config schema.
func (r *RuleSet) ConfigSchema() *hclext.BodySchema {
	r.config = config2.New()

	return hclext.ImpliedBodySchema(r.config)
}

// ApplyConfig applies the configuration to the ruleset.
func (r *RuleSet) ApplyConfig(body *hclext.BodyContent) error {
	predefinedResources := r.config.Resources
	predefinedModules := r.config.Modules
	r.config.Resources = make([]*config2.Resource, 0)
	r.config.Modules = make([]*config2.Resource, 0)

	diags := hclext.DecodeBody(body, nil, r.config)
	if diags.HasErrors() {
		return diags
	}

	r.config.Resources = append(predefinedResources, r.config.Resources...)
	r.config.Modules = append(predefinedModules, r.config.Modules...)

	return nil
}

// NewRunner creates a custom runner with the provided config.
func (r *RuleSet) NewRunner(runner tflint.Runner) (tflint.Runner, error) {
	return NewRunner(runner, r.config)
}
