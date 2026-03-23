# Test: Multi-liner key sorting in objects (two multi-line values out of order)
resource "example" "this" {
  tags = {
    "zzz" = jsonencode({
      foo = "bar"
    })

    "aaa" = jsonencode({
      baz = "qux"
    })
  }
}
