# Test: passing - module with for_each correctly as top-most attribute
module "example" {
  for_each = var.items
  source   = "./modules/example"

  name = each.value.name
}
