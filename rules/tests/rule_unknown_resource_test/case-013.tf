# Test: `data "X"` config does NOT cover `resource "X"` in code (no resource→data fallback)
resource "custom_lookup" "this" {
  id = "abc"
}
