# Test: failing - lifecycle at top instead of end of block
resource "aws_instance" "example" {
  lifecycle {
    create_before_destroy = true
  }

  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  iam_instance_profile = aws_iam_instance_profile.example

  depends_on = [
    aws_iam_role_policy.example
  ]
}
