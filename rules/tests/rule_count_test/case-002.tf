# Test: failing - count is not the top-most attribute
resource "random_password" "this" {
  length           = 16
  override_special = "!#$%&*()-_=+[]{}<>:?"
  special          = true

  count = 10
}
