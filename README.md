## s3put
Just put a file to S3, with single small binary tool

### IAM
s3put loads credentials via standard AWS methods such as `$HOME/.aws` and `AWS_XXX` environment variables.

s3put only requires `s3:PutObject` permission.

### Usage
To put a file `./foo` to `s3://your-bucket-name/foo`,
```bash
cat ./foo | s3put s3://your-bucket-name/foo
# or
s3put s3://your-bucket-name/foo < ./foo
```
