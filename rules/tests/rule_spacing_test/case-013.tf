# Test: redundant blank lines inside a map literal (before, between, after)
resource "example" "this" {
  tags = {

    "single_a" = "a"


    "single_b" = "b"

  }
}
