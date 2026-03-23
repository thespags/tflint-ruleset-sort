# TFLint Ruleset Sheldon

<div style=”text-align: right”><em>„This is my spot!”</em></div>

<br>

TFLint ruleset plugin that enforces consistent Terraform formatting — sorting, spacing, and structural conventions.

## Rules

| Rule | Description |
|------|-------------|
| [`sheldon_count`](docs/sheldon_count.md) | `count` placed at top of resource |
| [`sheldon_depends_on`](docs/sheldon_depends_on.md) | `depends_on` placed at end of resource |
| [`sheldon_for_each`](docs/sheldon_for_each.md) | `for_each` placed at top of resource |
| [`sheldon_key_attributes`](docs/sheldon_key_attributes.md) | Key attributes defined first and in priority order |
| [`sheldon_lifecycle`](docs/sheldon_lifecycle.md) | `lifecycle` block placed at end of resource |
| [`sheldon_sorting`](docs/sheldon_sorting.md) | Alphabetical sorting of blocks and dictionary keys |
| [`sheldon_source`](docs/sheldon_source.md) | `source` placed at top of module |
| [`sheldon_spacing`](docs/sheldon_spacing.md) | Consistent blank lines between attributes and blocks |

## Installation

You can install the plugin with `tflint --init`. Declare a config in
`.tflint.hcl` as follows:

```hcl
plugin "sheldon" {
  enabled = true

  version = "0.0.6"
  source  = "github.com/0x416e746f6e/tflint-ruleset-sheldon"
}
```

## Building the plugin

Clone the repository locally and run the following command:

```bash
make
```

You can easily install the built plugin with the following:

```bash
make install
```

You can run the built plugin like the following:

```bash
cat << EOF > .tflint.hcl
config {
  plugin_dir = "~/.tflint.d/plugins"
}

plugin "sheldon" {
  enabled = true
}
EOF

tflint
```

Some resources come with their key-attributes pre-defined.  However,
their set is far from being exhaustive. To define the
key-attributes for a resource/data blocks (so that `key_attributes`
rule picks them up) add them to in `.tflint.hcl` like follows:

```hcl
plugin "sheldon" {
  enabled = true

  resource "kubernetes_deployment" {
    key_attributes = ["metadata.namespace", "metadata.name"]
  }
}
```
