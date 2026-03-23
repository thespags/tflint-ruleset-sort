# Test: No resources to sort, passes with no issues
resource "single_attr" "example" {
  ami = "ami-a1b2c3d4"
}

module "non_resource" {
  source = "./modules/example"
}
