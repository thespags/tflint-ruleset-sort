# Test: missing blank line after multi-line value before single-line entry
resource "example" "this" {
  tags = {
    "multi_a" = jsonencode({
      foo = "bar"
    })
    "single" = "val"
  }
}
