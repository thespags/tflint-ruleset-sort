# Test: provider after for_each but not immediately — fails
resource "example" "this" {
  for_each = var.items

  name = "foo"

  provider = aws.west
}
