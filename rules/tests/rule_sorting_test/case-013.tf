# Test: Preprocessing drops trailing depends_on from sorting
resource "google_storage_bucket" "this" {
  project = "my-project"
  name    = "my-bucket"

  zzz = "last"
  aaa = "first"

  depends_on = [null_resource.this]
}
