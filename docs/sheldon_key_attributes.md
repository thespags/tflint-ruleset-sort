# Rule `key_attributes`

Makes sure that the key-attributes (those that uniquely identify a resource) are
placed on the top of the definition in the prioritized order.

## Example

```hcl
resource "kubernetes_service_account" "this" {
  metadata {
    name      = "service-account"
    namespace = kubernetes_namespace.this.metadata[0].name
  }
}
```

```text
Error: higher-priority key-attribute `namespace` should be defined before `name` (sheldon_key_attributes)

  on template.tf line 4:
   4:     namespace = kubernetes_namespace.this.metadata[0].name
```

## Module Example

Key attributes can also be configured for modules, keyed by their `source` value.

```hcl
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  other = "value"
  name  = "my-vpc"
  cidr  = "10.0.0.0/16"
}
```

```text
Error: key-attribute `name` should be defined before non-key `other` (sheldon_key_attributes)

  on template.tf line 5:
   5:   name  = "my-vpc"
```
