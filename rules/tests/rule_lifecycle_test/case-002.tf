# Test: failing - lifecycle is not at end (attributes follow it)
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  depends_on = [
    aws_iam_role_policy.example
  ]

  lifecycle {
    create_before_destroy = true
  }

  iam_instance_profile = aws_iam_instance_profile.example
}