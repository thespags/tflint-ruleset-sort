# Test: Object sorting inside tuple/list element (checkTupleConsExpr)
resource "example" "this" {
  list = [
    {
      "bbb" = "val"
      "aaa" = "val"
    },
  ]
}
