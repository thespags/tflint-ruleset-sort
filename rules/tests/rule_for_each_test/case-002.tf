# Test: failing - for_each is not the top-most attribute
resource "google_container_registry" "this" {
  project  = "foo"
  location = each.key
  for_each = ["a", "b"]
}
