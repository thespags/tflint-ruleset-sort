data "example_users" "this" {
  for_each = toset([
    "gamma",
    "alpha",
    "beta",
  ])

  username = each.value
}
