# Backend

## API
- Request: `GET /`
  - Response: `"Hello World"`
- Request: `GET /ws`
  - Prerequisite: authentication cookie
  - Effect: creates a new WebSocket connection

### User
- Request: `GET /api/user/`
  - Effect: creates a new user and sets authentication cookie unless already authenticated
  - Response: `{UserID: 1234, ...}`
- Request: `DELETE /api/user/`
  - Prerequisite: authentication cookie
  - Effect: deletes user

#### User Calendar
- Request: `PATCH /api/user/1234/availability/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Prerequisite: authentication cookie
  - Effect: Create new scheduled availability
- Request: `DELETE /api/group/1234/activity/5678/`
  - Prerequisite: authentication cookie
  - Effect: Delete scheduled availability by ID

### Group
- Request: `PATCH /api/group/ {Name: "abc"}`
  - Prerequisite: authentication cookie
  - Effect: creates a new group
  - Response: `{GroupID: 1234, ...}`
- Request: `GET /api/group/1234/`
  - Prerequisite: authentication cookie of user in group `1234`
  - Response: `{GroupID: 1234, Poll: {...}, Calendar: {...}, ...}`

#### Group Chat
- Request: `GET /api/group/1234/chat/?startTime=123456789` gets group chat messages starting at a Unix millisecond time
  - Prerequisite: authentication cookie of user in group `1234`
  - Response: `[{Sender: 5678, Timestamp: 123456789, Message: "hello", ...}, {...}, {...}]`
- Request: `PATCH /api/group/1234/chat/ {Message: "hello"}`
  - Prerequisite: authentication cookie of user in group `1234`
  - Effect: Send a chat message in group `1234`

#### Group Poll
- Request: `PUT /api/group/1234/poll/ {Title: "abc?", Options: ["a", "b", "c"]}`
  - Prerequisite: authentication cookie of user in group `1234`
  - Effect: Create/replace poll in group `1234`
- Request: `DELETE /api/group/1234/poll/`
  - Prerequisite: authentication cookie of user in group `1234`
  - Effect: Delete poll in group `1234`

#### Group Calendar
- Request: `PATCH /api/group/1234/activity/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Prerequisite: authentication cookie of user in group `1234`
  - Effect: Create new scheduled activity
- Request: `DELETE /api/group/1234/activity/5678/`
  - Prerequisite: authentication cookie of user in group `1234`
  - Effect: Delete scheduled activity by ID

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