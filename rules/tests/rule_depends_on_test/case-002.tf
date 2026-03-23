# Test: failing - depends_on is not the last attribute
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  depends_on = [
    aws_iam_role_policy.example
  ]

  iam_instance_profile = aws_iam_instance_profile.example
}
