# Test: provider after count — passes
resource "example" "this" {
  count = 3

  provider = aws.west

  name = "foo"
}
