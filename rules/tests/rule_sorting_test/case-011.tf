# Test: Object sorting inside parenthesized expression (checkParenthesesExpr)
resource "example" "this" {
  value = ({
    "bbb" = "val"
    "aaa" = "val"
  })
}
