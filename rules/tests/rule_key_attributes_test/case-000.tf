# Test: passing - unknown resource types are ignored
resource "unknown_type" "example" {
  ami = "ami-a1b2c3d4"
}

module "non_resource" {
  source = "./modules/example"
}
