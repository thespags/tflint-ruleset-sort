# Test: missing blank line after count and trailing blank before closing brace
resource "kubernetes_config_map" "this" {
  count = var.create_map ? 1 : 0
  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "config-map"
  }

  data = {
    "foo" = "bar"
  }

}

resource "aws_iam_user" "the-accounts" {
  for_each = toset(["Todd", "James", "Alice", "Dottie"])
  name     = each.key
}
