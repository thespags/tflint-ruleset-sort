# Rule `sort_provider`

Makes sure that `provider` meta-attribute is placed after `for_each`/`count`
but before all other attributes in `resource` or `data` block.

## Example

```hcl
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
  provider      = aws.west
}
```

```text
Error: `provider` must follow `for_each`/`count` (or be the top-most attribute) (sort_provider)

  on template.tf line 4:
   4:   provider = aws.west
```
