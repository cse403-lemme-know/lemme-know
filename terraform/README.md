# Terraform

## Prerequisites

1. Install `terraform`
2. Install `aws` CLI
3. Add AWS credentials to `~/.aws/credentials`
```toml
[cse403]
aws_access_key_id = {your AWS access key id}
aws_secret_access_key = {your AWS secret access key}
```
4. Create a [GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) with read access on `lemmeknow` (or switch to a different state provider).
5. Initialize Terraform with
```sh
terraform init -backend-config="password={your-github-token}"
```

## Provisioning infrastructure

To provision or update infrastructure, run the apply command:
```sh
terraform apply
```