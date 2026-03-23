# Test: passing - resources without depends_on are ignored
resource "no_depends_on" "example" {
  ami = "ami-a1b2c3d4"
}

module "non_resource" {
  source = "./modules/example"
}
