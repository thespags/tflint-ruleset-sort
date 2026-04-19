# Rule `sort_source`

Makes sure that `source` meta-attribute is placed at the top of `module` block,
or right after `for_each`/`count` when present.

## Example

```hcl
module "website_s3_bucket" {
  bucket_name = "<UNIQUE BUCKET NAME>"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }

  source = "./modules/aws-s3-static-website-bucket"
}
```

```text
Error: `source` must follow `for_each`/`count` (or be the top-most attribute) (sort_source)

  on template.tf line 9:
   9:   source = "./modules/aws-s3-static-website-bucket"
```
