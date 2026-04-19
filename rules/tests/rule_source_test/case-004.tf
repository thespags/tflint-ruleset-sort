# Test: passing - source after for_each
module "example" {
  for_each = var.items
  source   = "./modules/example"

  name = each.value.name
}
