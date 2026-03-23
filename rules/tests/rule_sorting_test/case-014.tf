# Test: Preprocessing drops trailing lifecycle block from sorting
resource "google_storage_bucket" "this" {
  project = "my-project"
  name    = "my-bucket"

  zzz = "last"
  aaa = "first"

  lifecycle {
    ignore_changes = [tags]
  }
}
