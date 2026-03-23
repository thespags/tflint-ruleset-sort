# Test: two unknown resources of same type - only first gets issue (dedup)
resource "beta_gadget" "first" {
  name = "one"
}

resource "beta_gadget" "second" {
  name = "two"
}
