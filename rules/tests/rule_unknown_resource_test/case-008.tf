# Test: mix of known and unknown - only unknown gets issue
resource "google_storage_bucket" "known" {
  project = var.project
  name    = "bucket"
}

resource "unknown_special" "this" {
  foo = "bar"
}

data "google_storage_bucket" "known_data" {
  project = var.project
  name    = "bucket"
}
