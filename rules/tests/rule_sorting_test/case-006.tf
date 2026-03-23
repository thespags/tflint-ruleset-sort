# Test: Single-line attribute precedes multi-line block; attribute precedes block
resource "example" "this" {
  nested {
    zzz_block {
      x = 1
      y = 2
    }

    aaa = "single"
  }
}
