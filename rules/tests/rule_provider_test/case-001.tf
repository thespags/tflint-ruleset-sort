# Test: provider already at top — passes
resource "example" "this" {
  provider = aws.west

  name = "foo"
}
