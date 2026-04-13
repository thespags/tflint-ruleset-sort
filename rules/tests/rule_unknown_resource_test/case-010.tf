# Test: duplicate unknown modules with same source - only first reported
module "foo_a" {
  source = "./modules/foo"
}

module "foo_b" {
  source = "./modules/foo"
}
