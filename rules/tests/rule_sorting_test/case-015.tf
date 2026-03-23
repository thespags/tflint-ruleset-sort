# Test: Key block dropping and inspection (kubernetes_service_account metadata)
resource "kubernetes_service_account" "this" {
  metadata {
    namespace = "default"
    name      = "my-sa"

    zzz = "last"
    aaa = "first"
  }
}
