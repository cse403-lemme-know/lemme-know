# Backend

## Prerequisites

1. Install `go` v1.21 or higher.
2. Install `make`.

## Running local backend

1. Run local backend:
```sh
make
```
2. Run the development frontend:
```sh
cd ../frontend
npm run build && npm run preview
```
2. Open http://localhost:8080.

## Deploying backend to AWS lambda

1. Apply Terraform to get credentials.
2. Deploy serverless backend:
```sh
make deploy
```

## API
- Open WebSocket: `/ws/`
  - Precondition: authentication cookie
  - Effect: creates a new WebSocket connection
  - Messages:
    - `{group: {groupId: 1234}}`
      - Meaning: means the group was updated and should be redownloaded.
    - `{message: {groupId: 1234, timestamp: 123456789, sender: 5678, content: "hello", ...}`
      - Meaning: delivers a chat message.
    - `{user: {userId: 5678, name: "Alex", status: "online" | "busy" | "offline"}}`
      - Meaning: a user profile changed.

### User
- Request: `GET /api/user/`
  - Effect: creates a new user and sets authentication cookie unless already authenticated.
  - Response: `{userId: 1234, name: "Alex", groups: [1234], ...}`
- Request: `PATCH /api/user/ {name: "Alex", status: "online" | "busy" | "offline"}`
  - Precondition: authentication cookie.
  - Effect: overwrites whichever profile settings were sent in the object.
- Request: `DELETE /api/user/`
  - Precondition: authentication cookie.
  - Effect: deletes user.

### Group
- Request: `GET /api/group/1234/`
  - Precondition: authentication cookie of user in group `1234`
  - Response: `{poll: {options: [{"a": [1234], ..}]}, availabilities: [{availabilityId: 5678, UserId: 5678, date: "9999-12-31", start: "8:00", end: "11:00"}], activities: [{activityId: 5678, Title: "abc", date: "9999-12-31", start: "9:00", end: "10:30", confirmed: [5678]}, ...], ...}`
- Request: `PATCH /api/group/1234/ {name: "Best Friends", calendarMode: "date" | "dayOfWeek"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: updates any group `1234` setting(s) passed in object.
  - Response: `{groupId: 1234}`.
- Request: `PATCH /api/group/ {name: "Friends"}`
  - Precondition: authentication cookie.
  - Effect: creates a new group with the specified name.
  - Response: `{groupId: 1234}`.

#### Chat
- Request: `GET /api/group/1234/chat/?Start=123456789&End=123456789` gets group chat messages starting at a Unix millisecond time (inclusive) and ending at a Unix millisecond time (inclusive).
  - Precondition: authentication cookie of user in group `1234`.
  - Response: `{Messages: [{sender: 5678, timestamp: 123456789, content: "hello", ...}, {...}], continue: true}`
    - Note: If `continue` is true, should request again (with higher start time) for more messages.
- Request: `PATCH /api/group/1234/chat/ {content: "hello"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Send a chat message in group `1234`.

#### Poll
- Request: `PUT /api/group/1234/poll/ {title: "abc?", options: ["a", "b", "c"]}`
  - Precondition: authentication cookie of user in group `1234`
  - Effect: Create/replace poll in group `1234`.
- Request: `PATCH /api/group/1234/poll/ {votes: ["b"]}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: cast vote(s) for option(s) in poll in group `1234`, replacing earlier vote(s)
- Request: `DELETE /api/group/1234/poll/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Dismiss poll in group `1234` to chat (immutable).

#### Availability
- Request: `PATCH /api/group/1234/availability/ {date: "9999-12-31", start: "9:00", end: "10:30"}`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Create new scheduled availability in group.
- Request: `DELETE /api/group/1234/availability/5678/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Delete scheduled availability by ID.

#### Activity
- Request: `PATCH /api/group/1234/activity/ {title: "abc", date: "9999-12-31", start: "9:00", end: "10:30"}`
  - Precondition: authentication cookie of user in group `1234`, activity doesn't overlap with others.
  - Effect: Create new scheduled activity.
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Delete scheduled activity by ID.