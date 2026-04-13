resource "google_storage_bucket" "this" {
  project = "my-project"
  name    = "my-bucket"

  location = "US"
}
