resource "example_project" "this" {
  enabled = true
  project = module.example.id
  member_ids = [
    99999,
    coalesce(var.x, 0),
    11111,
  ]
}
