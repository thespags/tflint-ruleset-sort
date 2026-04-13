# Test: module with dynamic source - no crash, reports empty source
module "dynamic" {
  source = var.module_source
}
