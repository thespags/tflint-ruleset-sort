# Test: Block alphabetical sorting (unsorted blocks at level > 0)
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.30.0"
    }
  }

  backend "gcs" {
    bucket                      = "terraform"
    impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
  }
}

provider "google" {
  impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
}
