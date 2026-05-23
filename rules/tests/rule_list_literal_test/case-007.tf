resource "example_service" "this" {
  args    = ["--flag", "value", "--other"]
  command = ["sh", "-c", "echo hello"]
}
