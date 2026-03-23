# Test: Multi-liner alphabetical sorting (two multi-line blocks out of order)
resource "example" "this" {
  nested {
    zzz_block {
      a = 1
    }

    aaa_block {
      b = 2
    }
  }
}
