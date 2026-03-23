# Test: Attribute after dynamic block (single-line precedes multi-line)
resource "example" "this" {
  nested {
    dynamic "zzz_block" {
      content {
        x = 1
      }
    }

    aaa = "val"
  }
}
