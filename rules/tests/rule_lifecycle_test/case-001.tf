# Test: passing - lifecycle is correctly before depends_on
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  iam_instance_profile = aws_iam_instance_profile.example

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_iam_role_policy.example
  ]
}
