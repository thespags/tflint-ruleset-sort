# Test: two different unknown modules - each gets an issue
module "zeta_mod_a" {
  source = "./modules/a"
}

module "zeta_mod_b" {
  source = "./modules/b"
}
