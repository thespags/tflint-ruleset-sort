locals {
  data = {
    for x in var.items : x => {
      aaa = "val"
      nested = {
        zzz = "val"
        yyy = "val"
      }
    }
  }
}
