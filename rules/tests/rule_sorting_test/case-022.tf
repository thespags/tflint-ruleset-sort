locals {
  data = {
    for x in var.items : x => {
      ccc = "val3"
      bbb = "val2"
      aaa = "val1"

      nested = {
        zzz = 1
        yyy = 2
      }
    }
  }
}
