# Rule `sort_unknown_resource`

Warns when a resource, data source, or module is encountered that is not
defined in the plugin configuration. This helps identify resources that may
need key-attributes configured.

`resource "X"` and `data "X"` are looked up independently with a one-way
fallback: a `data "X"` block in code is considered "known" if either a
`data "X"` or a `resource "X"` block is declared in `.tflint.hcl`. A
`resource "X"` block is considered "known" only if a `resource "X"` block is
declared — it does **not** fall back to a `data "X"` config.

## Example

```text
Warning: unknown resource `aws_instance` — consider adding it to the sort plugin configuration (sort_unknown_resource)
```
