package rules

import (
	"fmt"
	"sync"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/thespags/tflint-ruleset-sort/custom"
	"github.com/thespags/tflint-ruleset-sort/project"
	"github.com/thespags/tflint-ruleset-sort/visit"
)

// UnknownResourceRule warns if the linter encounters unknown resource.
type UnknownResourceRule struct {
	tflint.DefaultRule

	mu            *sync.Mutex
	seenResources map[string]struct{}
	seenModules   map[string]struct{}
}

// NewUnknownResourceRule creates a new UnknownResourceRule.
func NewUnknownResourceRule() *UnknownResourceRule {
	return &UnknownResourceRule{
		mu:            &sync.Mutex{},
		seenResources: map[string]struct{}{},
		seenModules:   map[string]struct{}{},
	}
}

// Name returns the name of the rule.
func (*UnknownResourceRule) Name() string {
	return project.RuleName("unknown_resource")
}

// Enabled returns whether the rule is enabled by default.
func (*UnknownResourceRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (*UnknownResourceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the reference link for the rule.
func (r *UnknownResourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the key-attributes (those that uniquely identify the
// resource) are put on top of the resource definition.
func (r *UnknownResourceRule) Check(rr tflint.Runner) error {
	switch runner := rr.(type) {
	case *custom.Runner:
		return visit.Blocks(runner, func(block *hclsyntax.Block, _ []byte) error {
			if block.Type == "resource" || block.Type == "data" {
				return r.checkBlock(runner, runner.Resources, block)
			}

			if block.Type == "module" {
				return r.checkBlock(runner, runner.Modules, block)
			}

			return nil
		})

	default:
		return nil
	}
}

func (r *UnknownResourceRule) checkBlock(
	runner *custom.Runner,
	resources map[string]*custom.Resource,
	block *hclsyntax.Block,
) error {
	kind := block.Labels[0]
	if _, kindIsKnown := resources[kind]; kindIsKnown {
		return nil
	}

	r.mu.Lock()

	_, seen := r.seenResources[kind]
	if !seen {
		r.seenResources[kind] = struct{}{}
	}
	r.mu.Unlock()

	if !seen {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("key-attributes for resource type `%s` are not configured", kind),
			block.LabelRanges[0],
		)
	}

	return nil
}
