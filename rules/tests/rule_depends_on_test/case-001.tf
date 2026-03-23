# Test: passing - depends_on is correctly the last attribute
data "aws_ami" "example" {
  most_recent = true
  owners      = ["self"]
  depends_on  = [aws_vpc.example]
}
