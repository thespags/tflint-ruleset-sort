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
	seenData      map[string]struct{}
	seenResources map[string]struct{}
	seenModules   map[string]struct{}
}

// NewUnknownResourceRule creates a new UnknownResourceRule.
func NewUnknownResourceRule() *UnknownResourceRule {
	return &UnknownResourceRule{
		mu:            &sync.Mutex{},
		seenData:      map[string]struct{}{},
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
			resource := runner.Lookup(block)
			if resource != nil {
				return nil
			}

			return r.checkIssue(runner, block)
		})

	default:
		return nil
	}
}

func (r *UnknownResourceRule) checkIssue(runner *custom.Runner, block *hclsyntax.Block) error {
	var (
		seenMap map[string]struct{}
		kind    string
	)

	switch block.Type {
	case "data":
		kind = block.Labels[0]
		seenMap = r.seenData
	case "resource":
		kind = block.Labels[0]
		seenMap = r.seenResources
	case "module":
		kind = custom.GetSource(block)
		seenMap = r.seenModules
	}

	if seenMap == nil {
		// not a resource we care about
		return nil
	}

	r.mu.Lock()

	_, seen := seenMap[kind]
	if !seen {
		seenMap[kind] = struct{}{}
	}
	r.mu.Unlock()

	if seen {
		return nil
	}

	return runner.EmitIssue(
		r,
		fmt.Sprintf("key-attributes for %s type `%s` are not configured", block.Type, kind),
		block.LabelRanges[0],
	)
}
