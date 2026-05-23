# Rule `sort_list_literal`

Sorts the elements of list literals (HCL tuples) everywhere they appear:
top-level attribute values, function-call arguments, nested object values,
`for` expression results, etc. Comparison is numeric-aware (so `"v2" < "v10"`
and `3 < 100`). Trailing inline comments (`// ...` or `# ...`) on each
element travel with the element when the rule auto-fixes. Contiguous
single-line comment lines (`// ...` or `# ...`) directly above an element
are treated as leading comments for that element and travel with it as well.

## Example

```hcl
data "example_users" "this" {
  for_each = toset([
    "gamma",
    "alpha",
    "beta",
  ])

  username = each.value
}
```

```text
Error: list-literal is not sorted: expected `alpha` at position 0, got `gamma` (sort_list_literal)

  on template.tf line 3:
   3:     "gamma",
```

## Supported element kinds

- String literals (`"foo"`)
- Number and boolean literals (`42`, `true`)
- Variable / data / module / local references (`var.x`, `data.y.z`,
  `module.a.b`, `local.c`) — sorted by their full dotted path

If any element of the list is something the rule can't lexically order — a
function call, a complex template with interpolations, an arithmetic
expression, etc. — the rule silently skips that list.

## Opting out per attribute (`skip_sort_literals`)

Some attributes have semantically meaningful list order: `command`, `args`,
URL `path` matchers, and so on. Add their names to the resource's
`skip_sort_literals` config to leave their values alone. The opt-out applies
to the attribute's entire expression subtree, so nested lists inside a
skipped attribute are also left untouched.

```hcl
plugin "sort" {
  enabled = true

  resource "example_service" {
    skip_sort_literals = ["args", "command"]
  }
}
```

Lookup is per-block-type. A `data "X"` block falls back to a `resource "X"`
config if no `data "X"` is declared, matching the existing `key_attributes`
semantics.

## Disabling the rule

To turn list-literal sorting off entirely, disable the rule in `.tflint.hcl`:

```hcl
plugin "sort" {
  enabled = true
}

rule "sort_list_literal" {
  enabled = false
}
```
