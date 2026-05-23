resource "example_project" "this" {
  enabled = true
  project = module.example.id
  member_ids = [
    99999, // zeta
    11111, // alpha
    55555, // mu
  ]
}
