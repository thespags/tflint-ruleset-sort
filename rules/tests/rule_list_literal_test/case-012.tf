resource "example_project" "this" {
  target_project_ids = [
    // TODO: foobar
    // see https://example.com/issues/42
    module.bravo.id,
    module.alpha.id,
    module.charlie.id,
  ]
}
