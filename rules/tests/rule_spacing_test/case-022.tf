resource "google_storage_bucket" "this" {
  count = var.create ? 1 : 0

  project = "my-project"
  name    = "my-bucket"

  location = "US"
}
