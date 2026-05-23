locals {
  result = { for x in var.items : x.name => [x.zone, x.region] }
}
