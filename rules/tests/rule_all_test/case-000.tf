module "crusoeenergy_stack_validate" {
  source = "./modules/project"

  merge_method  = "merge"
  namespace_id  = module.crusoeenergy.id
  path          = "stack-validate"
  squash_option = "default_off"
  branch_protection = {
    code_owner_approval_required = false
  }
}
