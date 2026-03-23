# Test: missing blank line between multi-line block and single-line attr
resource "example" "this" {
  metadata {
    name = "foo"
  }
  attr_a = "a"
}
