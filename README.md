![GitHub tag](https://img.shields.io/github/v/tag/thespags/tflint-ruleset-sort)
![Build](https://img.shields.io/github/actions/workflow/status/thespags/tflint-ruleset-sort/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/thespags/tflint-ruleset-sort)
![License](https://img.shields.io/github/license/thespags/tflint-ruleset-sort)
[![Go Report Card](https://goreportcard.com/badge/github.com/thespags/tflint-ruleset-sort)](https://goreportcard.com/report/github.com/thespags/tflint-ruleset-sort)
[![codecov](https://codecov.io/gh/thespags/tflint-ruleset-sort/branch/main/graph/badge.svg)](https://codecov.io/gh/thespags/tflint-ruleset-sort)

# TFLint Ruleset Sort

TFLint ruleset plugin that enforces consistent Terraform formatting — sorting, spacing, and structural conventions.

> **Note:** This project is a fork of [tflint-ruleset-sheldon](https://github.com/0x416e746f6e/tflint-ruleset-sheldon) by [@0x416e746f6e](https://github.com/0x416e746f6e), which was published without a license. This fork is maintained independently.

## Rules

| Rule | Description |
|------|-------------|
| [`sort_count`](docs/sort_count.md) | `count` placed at top of resource or module |
| [`sort_depends_on`](docs/sort_depends_on.md) | `depends_on` placed at end of resource |
| [`sort_for_each`](docs/sort_for_each.md) | `for_each` placed at top of resource or module |
| [`sort_key_attributes`](docs/sort_key_attributes.md) | Key attributes defined first and in priority order |
| [`sort_lifecycle`](docs/sort_lifecycle.md) | `lifecycle` block placed at end of resource |
| [`sort_provider`](docs/sort_provider.md) | `provider` placed after `for_each`/`count` |
| [`sort_sorting`](docs/sort_sorting.md) | Alphabetical sorting of blocks and dictionary keys |
| [`sort_source`](docs/sort_source.md) | `source` placed after `for_each`/`count` in module |
| [`sort_spacing`](docs/sort_spacing.md) | Consistent blank lines between attributes and blocks |
| [`sort_unknown_resource`](docs/sort_unknown_resource.md) | Warns on resources not in the configuration |

## Installation

You can install the plugin with `tflint --init`. Declare a config in
`.tflint.hcl` as follows:

```hcl
plugin "sort" {
  enabled = true

  version = "0.0.5"
  source  = "github.com/thespags/tflint-ruleset-sort"
}
```

## Building the plugin`

Clone the repository locally and run the following command:

With mise,
```bash
mise install
```

Build the plugin with:

```bash
go build ./...
```

You can install the built plugin with the following:

```bash
mise run install
```

You can run the built plugin like the following:

```bash
cat << EOF > .tflint.hcl
config {
  plugin_dir = "~/.tflint.d/plugins"
}

plugin "sort" {
  enabled = true
}
EOF

tflint
```

Some resources come with their key-attributes pre-defined.  However,
their set is far from being exhaustive. To define the
key-attributes for resource/data/module blocks (so that `key_attributes`
rule picks them up) add them to `.tflint.hcl` like follows:

```hcl
plugin "sort" {
  enabled = true

  resource "kubernetes_deployment" {
    key_attributes = ["metadata.namespace", "metadata.name"]
  }

  module "terraform-aws-modules/vpc/aws" {
    key_attributes = ["name", "cidr"]
  }
}
```

Module blocks are keyed by their `source` value rather than the module name,
since module names are user-defined and the source uniquely identifies the module.
