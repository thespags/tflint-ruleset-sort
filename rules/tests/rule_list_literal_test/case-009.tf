resource "example" "this" {
  items = [
    "${var.z}",
    "${var.a}",
    "${var.m}",
  ]
}
