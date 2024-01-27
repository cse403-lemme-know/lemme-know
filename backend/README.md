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

## API
- Request: `GET /`
  - Response: `"Hello World"`
- Request: `GET /ws`
  - Precondition: authentication cookie
  - Effect: creates a new WebSocket connection

### User
- Request: `GET /api/user/`
  - Effect: creates a new user and sets authentication cookie unless already authenticated
  - Response: `{UserID: 1234, Name: "Alex", GroupIDs: [1234], ...}`
- Request: `DELETE /api/user/`
  - Precondition: authentication cookie
  - Effect: deletes user

#### Calendar
- Request: `PATCH /api/user/1234/availability/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Precondition: authentication cookie
  - Effect: Create new scheduled availability
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie
  - Effect: Delete scheduled availability by ID

### Group
- Request: `GET /api/group/1234/`
  - Precondition: authentication cookie of user in group `1234`
  - Response: `{GroupID: 1234, Poll: {Options: [{"a": [1234], ..}]}, Activites: [{ActivityID: 5678, Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}, ...], ...}`
- Request: `PATCH /api/group/ {Name: "Friends"}`
  - Precondition: authentication cookie
  - Effect: creates a new group with the specified name
  - Response: `{GroupID: 1234}`

#### Chat
- Request: `GET /api/group/1234/chat/?Start=123456789&End=123456789` gets group chat messages starting at a Unix millisecond time (inclusive) and ending at a Unix millisecond time (exclusive).
  - Precondition: authentication cookie of user in group `1234`
  - Response: `{Messages: [{Sender: 5678, Timestamp: 123456789, Message: "hello", ...}, {...}], End: 123456789}` (response.End may be less than request.End for pagination purposes>)
- Request: `PATCH /api/group/1234/chat/ {Message: "hello"}`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Send a chat message in group `1234`

#### Poll
- Request: `PUT /api/group/1234/poll/ {Title: "abc?", Options: ["a", "b", "c"]}`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Create/replace poll in group `1234`
- Request: `PATCH /api/group/1234/poll/ {Votes: ["b"]}`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: cast vote(s) for option(s) in poll in group `1234`, replacing earlier vote(s)
- Request: `DELETE /api/group/1234/poll/`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Dismiss poll in group `1234` to chat (immutable)

#### Calendar
- Request: `PATCH /api/group/1234/activity/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Precondition: authentication cookie of user in group `1234`, activity doesn't overlap with others
  - Effect: Create new scheduled activity
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Delete scheduled activity by ID