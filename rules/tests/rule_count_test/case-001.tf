# Test: passing - count is correctly the top-most attribute
resource "random_password" "this" {
  count = 10

  length           = 16
  override_special = "!#$%&*()-_=+[]{}<>:?"
  special          = true
}
