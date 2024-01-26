1. Add AWS credentials to `~/.aws/credentials`
```toml
[cse453]
aws_access_key_id = {your AWS access key id}
aws_secret_access_key = {your AWS secret access key}
```
2. Create a GitHub personal access token with read access on `lemmeknow`.
3. Initialize Terraform with
```sh
terraform init -backend-config="password={your-github-token}"
```