# Backend

See [`../README.md`](../README.md) for instructions.

## API
- Open WebSocket: `/ws/`
  - Precondition: Authentication cookie
  - Effect: Creates a new WebSocket connection
  - Messages:
    - `{group: {groupId: 1234}}`
      - Meaning: Means the group was updated and should be redownloaded.
    - `{message: {groupId: 1234, timestamp: 123456789, sender: 5678, content: "hello", ...}`
      - Meaning: Delivers a chat message.
    - `{user: {userId: 5678, name: "Alex", status: "online" | "busy" | "offline"}}`
      - Meaning: A user profile changed.

### Push
- Request: `GET /api/push/`
  - Precondition: Authentication cookie
  - Response: `{vapidPublicKey: "XYZW"}`
- Request: `PATCH /api/push/ {"endpoint": "https://example.com", "keys": {"auth": "XYZW", "p256dh": "XYZW"}}`
  - Precondition: Authentication cookie
  - Effect: Adds WebPush subscription.
  - Messages:
    - `{message: {groupId: 1234, timestamp: 123456789, sender: "Bob", content: "hello", ...}`
      - Meaning: Delivers a chat message.
    - `{reminder: {groupId: 1234, timestamp: 123456789, content: "hello", ...}`
      - Meaning: Delivers a reminder.

### User
- Request: `GET /api/user/`
  - Effect: Creates a new user and sets authentication cookie unless already authenticated.
  - Response: `{userId: 1234, name: "Alex", groups: [1234], ...}`
- Request: `PATCH /api/user/ {name: "Alex", status: "online" | "busy" | "offline"}`
  - Precondition: Authentication cookie.
  - Effect: overwrites whichever profile settings were sent in the object.
- Request: `DELETE /api/user/`
  - Precondition: Authentication cookie.
  - Effect: Deletes user.

### Group
- Request: `GET /api/group/1234/`
  - Precondition: Authentication cookie.
  - Effect: User `1234` joins the group if they weren't in it already.
  - Note: Most group API operations return nothing and instead issue an unsolicited notification for all participating clients to use this API to re-download the group.
  - Response: `{poll: {options: [{"a": [1234], ..}]}, availabilities: [{availabilityId: 5678, UserId: 5678, date: "9999-12-31", start: "8:00", end: "11:00"}], activities: [{activityId: 5678, Title: "abc", date: "9999-12-31", start: "9:00", end: "10:30", confirmed: [5678]}, ...], tasks: [{taskId: 2345, title: "prepare food & drinks", assignee: 5678, complete: true}, ...], ...}`
- Request: `PATCH /api/group/1234/ {name: "Best Friends", calendarMode: "date" | "dayOfWeek"}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Updates any group `1234` setting(s) passed in object.
  - Response: `{groupId: 1234}`.
- Requet: `DELETE /api/group/1234/`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Leaves the group, deleting it if last to leave.
- Request: `PATCH /api/group/ {name: "Friends"}`
  - Precondition: Authentication cookie.
  - Effect: Creates a new group with the specified name.
  - Response: `{groupId: 1234}`.

#### Activity
- Request: `PATCH /api/group/1234/activity/ {title: "abc", date: "9999-12-31", start: "9:00", end: "10:30"}`
  - Precondition: Authentication cookie of user in group `1234`, activity doesn't overlap with others.
  - Effect: Create new scheduled activity.
- Request: `PATCH /api/group/1234/activity/5678/ {title: "abc", date: "9999-12-31", start: "9:00", end: "10:30", confirm: true}`
  - Precondition: Authentication cookie of user in group `1234`, activity doesn't overlap with others.
  - Effect: Edit scheduled activity, notably by marking whether the requesting user confirms their attendance. If the date, start, or end is changed, the availability of others will be erased.
- Request: `DELETE /api/group/1234/activity/5678/`
  - Precondition: authentication cookie of user in group `1234`.
  - Effect: Delete scheduled activity by ID.

#### Availability
- Request: `PATCH /api/group/1234/availability/ {date: "9999-12-31", start: "9:00", end: "10:30"}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Create new scheduled availability in group.
- Request: `DELETE /api/group/1234/availability/5678/`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Delete scheduled availability by ID.

#### Chat
- Request: `GET /api/group/1234/chat/?start=123456789&end=123456789` gets group chat messages starting at a Unix millisecond time (inclusive) and ending at a Unix millisecond time (inclusive).
  - Precondition: Authentication cookie of user in group `1234`.
  - Response: `{messages: [{sender: 5678, timestamp: 123456789, content: "hello", ...}, {...}], continue: true}`
    - Note: If `continue` is true, should request again (with higher start time of `messages[last].timestamp + 1`) for more messages.
- Request: `PATCH /api/group/1234/chat/ {content: "hello"}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Send a chat message in group `1234`.

#### Poll
- Request: `PUT /api/group/1234/poll/ {title: "abc?", options: ["a", "b", "c"]}`
  - Precondition: Authentication cookie of user in group `1234`
  - Effect: Create/replace poll in group `1234`.
- Request: `PATCH /api/group/1234/poll/ {votes: ["b"]}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Cast vote(s) for option(s) in poll in group `1234`, replacing earlier vote(s)
- Request: `DELETE /api/group/1234/poll/`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Dismiss poll in group `1234` to chat (immutable).

#### Task
- Request: `PATCH /api/group/1234/task/ {title: "prepare food", assignee: 4567}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Create new task for self in group for a given user (default to self if `assignee` unspecified or unknown).
- Request: `PATCH /api/group/1234/task/5678/ {title: "prepare food & drinks", assignee: 5678, complete: true}`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Update title, assignee, and/or completion status.
- Request: `DELETE /api/group/1234/task/5678/`
  - Precondition: Authentication cookie of user in group `1234`.
  - Effect: Delete task by ID.