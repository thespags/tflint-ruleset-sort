# Test: two unknown data sources of same type - only first gets issue (dedup)
data "delta_query" "first" {
  id = "123"
}

data "delta_query" "second" {
  id = "456"
}
