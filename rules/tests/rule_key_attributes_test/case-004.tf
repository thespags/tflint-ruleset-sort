# Test: failing - namespace should be defined before name
resource "kubernetes_service_account" "this" {
  metadata {
    name      = "service-account"
    namespace = kubernetes_namespace.this.metadata[0].name
  }
}
