locals {
  data = {
    for x in var.items : x => {
      ccc = "val"
      bbb = "val"
      aaa = "val"
    }
  }
}
