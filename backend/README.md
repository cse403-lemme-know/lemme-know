# Backend

## Prerequisites

1. Install `go` v1.20 or higher
2. Install `make`

## Running local service

1. Run local service
```sh
make
```
2. Open http://localhost:8080

## Deploying to AWS lambda

1. Run Terraform to get credentials
2. Deploy serverless service
```sh
make deploy
```