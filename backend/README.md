# Backend

## Prerequisites

1. Install `go` v1.20 or higher.
2. Install `make`.

## Running local backend

1. Run local backend:
```sh
make
```
2. Run the development frontend (see [../frontend/README.md](../frontend/README.md))
2. Open http://localhost:8080.

## Deploying backend to AWS lambda

1. Apply Terraform to get credentials.
2. Deploy serverless backend:
```sh
make deploy
```

## API
- Request: `GET /`
  - Response: `"Hello World"`
- Open WebSocket: `/ws/`
  - Precondition: authentication cookie
  - Effect: creates a new WebSocket connection
  - Messages:
    - `{GroupID: 1234}`
      - Meaning: means the group was updated and should be redownloaded.
    - `{Message: {GroupID: 1234, Timestamp: 123456789, Sender: 5678, Message: "hello", ...}`
      - Meaning: delivers a chat message.
    - `{UserID: 5678, Name: "Alex", Status: "online" | "busy" | "offline"}`
      - Meaning: a user profile changed.

### User
- Request: `GET /api/user/`
  - Effect: creates a new user and sets authentication cookie unless already authenticated.
  - Response: `{UserID: 1234, Name: "Alex", GroupIDs: [1234], ...}`
- Request: `PATCH /api/user/ {Name: "Alex", Status: "online" | "busy" | "offline"}`
  - Precondition: authentication cookie.
  - Effect: overwrites whichever profile settings were sent in the object.
- Request: `DELETE /api/user/`
  - Precondition: authentication cookie.
  - Effect: deletes user.

### Group
- Request: `GET /api/group/1234/`
  - Precondition: authentication cookie of user in group `1234`
  - Response: `{GroupID: 1234, Poll: {Options: [{"a": [1234], ..}]}, Availabilities: [{AvailabilityID: 5678, UserID: 5678, Date: "9999-12-31", Start: "8:00", End: "11:00"}], Activities: [{ActivityID: 5678, Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}, ...], ...}`
- Request: `PATCH /api/group/1234/ {Name: "Best Friends", CalendarTitle: "Hikes", CalendarMode: "date" | "dayOfWeek"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: updates any group `1234` setting(s) passed in object.
  - Response: `{GroupID: 1234}`.
- Request: `PATCH /api/group/ {Name: "Friends"}`
  - Precondition: authentication cookie.
  - Effect: creates a new group with the specified name.
  - Response: `{GroupID: 1234}`.

#### Chat
- Request: `GET /api/group/1234/chat/?Start=123456789&End=123456789` gets group chat messages starting at a Unix millisecond time (inclusive) and ending at a Unix millisecond time (inclusive).
  - Precondition: authentication cookie of user in group `1234`.
  - Response: `{Messages: [{Sender: 5678, Timestamp: 123456789, Message: "hello", ...}, {...}], End: 123456789}`
    - Note: `response.End` may be less than `request.End` for pagination purposes.
- Request: `PATCH /api/group/1234/chat/ {Message: "hello"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Send a chat message in group `1234`.

#### Poll
- Request: `PUT /api/group/1234/poll/ {Title: "abc?", Options: ["a", "b", "c"]}`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Create/replace poll in group `1234`.
- Request: `PATCH /api/group/1234/poll/ {Votes: ["b"]}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: cast vote(s) for option(s) in poll in group `1234`, replacing earlier vote(s)
- Request: `DELETE /api/group/1234/poll/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Dismiss poll in group `1234` to chat (immutable).

#### Availability
- Request: `PATCH /api/group/1234/availability/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Create new scheduled activity.
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Delete scheduled availability by ID.

#### Activity
- Request: `PATCH /api/group/1234/activity/ {Title: "abc", Date: "9999-12-31", Start: "9:00", End: "10:30"}`
  - Precondition: authentication cookie of user in group `1234`, activity doesn't overlap with others.
  - Effect: Create new scheduled activity.
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Delete scheduled activity by ID.