# Test: known module (configured in config) - no issue
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
}
