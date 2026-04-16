# Test: provider not at expected position — fails
resource "example" "this" {
  name = "foo"

  provider = aws.west
}
