# Test: Single-line key/value should be placed before multi-line in object
resource "example" "this" {
  tags = {
    "multi" = jsonencode({
      foo = "bar"
    })
    "aaa" = "single"
  }
}
