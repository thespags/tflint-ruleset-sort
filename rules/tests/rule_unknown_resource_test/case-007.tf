# Test: non resource/data/module blocks - no issues
terraform {
  required_version = ">= 1.0"
}

provider "google" {
  project = var.project
}

locals {
  name = "example"
}
