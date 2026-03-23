# Test: passing - resources without for_each are ignored
resource "no_for_each" "example" {
  ami = "ami-a1b2c3d4"
}

module "non_resource" {
  source = "./modules/example"
}
