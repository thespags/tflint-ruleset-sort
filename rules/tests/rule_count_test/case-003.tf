# Test: passing - module with count correctly as top-most attribute
module "example" {
  count  = 3
  source = "./modules/example"

  name = "instance-${count.index}"
}
