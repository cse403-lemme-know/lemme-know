# Backend

## API
- `GET /`` returns a hello world message
- `GET /api/user/` creates a new user and stores cookie in browser if needed
- `PUT /api/group/` creates a new group
- `GET /api/group/1234/` gets a group by id
- `GET /api/group/1234/chat/?startTime=123456789` gets group chat messages starting at a Unix millisecond time
- `GET /ws` opens a new WebSocket connection

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