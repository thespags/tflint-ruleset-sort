# Test: Object key alphabetical sorting (unsorted keys in nested object)
resource "something" "this" {
  template {
    metadata {
      annotations = {
        "def" = "blah-blah"
        "abc" = "yada-yada"
      }
    }
  }
}
