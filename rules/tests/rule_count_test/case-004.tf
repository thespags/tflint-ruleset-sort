# Test: failing - module with count not as top-most attribute
module "example" {
  source = "./modules/example"
  name   = "foo"
  count  = 3
}
