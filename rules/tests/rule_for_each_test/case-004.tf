# Test: failing - module with for_each not as top-most attribute
module "example" {
  source   = "./modules/example"
  name     = "foo"
  for_each = var.items
}
