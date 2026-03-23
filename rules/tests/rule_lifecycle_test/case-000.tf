# Test: passing - resources without lifecycle are ignored
resource "no_lifecycle" "example" {
}

module "non_resource" {
}