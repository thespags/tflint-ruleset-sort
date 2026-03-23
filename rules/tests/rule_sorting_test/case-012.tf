# Test: Preprocessing drops leading count from sorting
resource "google_storage_bucket" "this" {
  count = 1

  project = "my-project"
  name    = "my-bucket"

  zzz = "last"
  aaa = "first"
}
