data "example_users" "this" {
  for_each = toset([
    "alpha",
    "beta",
    "gamma",
  ])

  username = each.value
}
