module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  other = "value"
  name  = "my-vpc"
  cidr  = "10.0.0.0/16"
}
