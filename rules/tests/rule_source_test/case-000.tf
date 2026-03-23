# Test: passing - non-module blocks are ignored
resource "not_a_module" "example" {
  ami = "ami-a1b2c3d4"
}

data "not_a_module" "example" {
  ami = "ami-a1b2c3d4"
}
