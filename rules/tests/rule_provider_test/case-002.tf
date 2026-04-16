# Test: provider after for_each — passes
resource "example" "this" {
  for_each = var.items

  provider = aws.west

  name = "foo"
}
