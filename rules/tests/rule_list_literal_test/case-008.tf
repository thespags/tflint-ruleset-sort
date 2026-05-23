resource "example" "this" {
  items = [
    "zeta",
    "prefix-${var.x}",
    "alpha",
  ]
}
