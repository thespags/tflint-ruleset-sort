# Test: failing - namespace before name in deeply nested metadata
resource "kubernetes_manifest" "qux" {
  manifest {
    apiVersion = "eek"
    kind       = "ook"

    metadata {
      name      = "foo"
      namespace = "bar"
    }
  }
}
