resource "example_acl" "this" {
  branch  = "main"
  project = module.example.id
  allowed_members = [
    var.foo.email,
    var.bar.email,
    "xyz",
  ]
}
