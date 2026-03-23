# Test: redundant blank lines inside function call arguments
resource "example" "this" {
  attr = concat(

    ["a"],


    ["b"],

  )
}
