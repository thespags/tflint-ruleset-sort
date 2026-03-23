# Test: Key block step-inside checks nested non-key content
resource "kubernetes_service_account" "this" {
  metadata {
    namespace = "default"
    name      = "my-sa"

    labels = {
      "bbb" = "val"
      "aaa" = "val"
    }
  }
}
