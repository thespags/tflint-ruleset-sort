# Test: empty resource and non-resource module (no issues expected)
resource "empty" "example" {
}

module "non_resource" {
  source = "./modules/example"
}
