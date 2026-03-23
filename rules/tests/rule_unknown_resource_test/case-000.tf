# Test: known resource type (no issue expected)
resource "google_storage_bucket" "this" {
  project = var.project
  name    = "my-bucket"
}
