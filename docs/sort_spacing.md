# Rule `sort_spacing`

Normalises blank-lines in the sources.

- No multiple consecutive blank-lines.
- No unnecessary single blank-lines.
- Comments are ignored.
- When key-attributes are configured, they are grouped together without blank lines,
  with a blank line separating them from the remaining attributes. This applies to
  both resources/data blocks and modules (where `source` is included in the group).

## Example

```hcl
resource "kubernetes_config_map" "this" {
  count = var.create_map ? 1 : 0
  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "config-map"
  }

  data = {
    "foo" = "bar"
  }

}
```

```text
Error: attribute `count` must be separated from the rest of the definition by an extra line (sort_spacing)

  on case-008.tf line 2:
   2:   count = var.create_map ? 1 : 0

Error: multi-line element must be separated from the previous one by an extra line (sort_spacing)

  on case-008.tf line 3:
   3:   metadata {
   4:     namespace = kubernetes_namespace.this.metadata[0].name
   5:     name      = "config-map"
   6:   }

Error: 1 redundant empty line in front (sort_spacing)

  on case-008.tf line 12:
  12: }
```

## Key-attribute grouping

When a resource or module has key-attributes configured, they are grouped with
no blank lines between them, and a blank line is required after the last
key-attribute before the remaining attributes.

```hcl
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "my-vpc"
  cidr   = "10.0.0.0/16"

  other = "value"
}
```

This behavior can be disabled by setting `separate_key_attributes` to `false`
in the `Disabled` configuration of the runner.
