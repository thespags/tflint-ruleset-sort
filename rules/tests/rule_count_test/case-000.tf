# Test: passing - resources without count are ignored
resource "no_count" "example" {
  ami = "ami-a1b2c3d4"
}

module "non_resource" {
  source = "./modules/example"
}
