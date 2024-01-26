1. Create a GitHub personal access token with read access on `lemmeknow`.
2. Initialize Terraform with
```sh
terraform init -backend-config="password={your-github-token}"
```