# Test: data block configured via `data "X"` block in .tflint.hcl - no issue
data "custom_lookup" "this" {
  id = "abc"
}
