package main

type Group struct {
	GroupID GroupID `dynamo:",hash"` // Hash key, a.k.a. partition key
	//Time      time.Time // Range key, a.k.a. sort key

	Name string
	//Count     int                  `dynamo:",omitempty"` // Omits if zero value
	Poll           *Poll
	Members        []UserID
	Activities     []Activity
	Availabilities []Availability
}

type User struct {
	UserID      UserID `dynamo:",hash"` // Hash key, a.k.a. partition key
	Name        string
	Groups      []GroupID           `dynamo:",set"`
	Connections []ConnectionID      `dynamo:",set"`
	Schedules   map[string]Schedule `dynamo:",set"`
}

type Poll struct {
	Title     string
	Timestamp uint64
	Options   []PollOption
	DoneFlag  bool
}

type Message struct {
	GroupID   GroupID `dynamo:",hash"`
	Timestamp uint64  `dynamo:",range"`
	Content   string
	Sender    UserID
}

type PollOption struct {
	Name  string
	Votes []UserID `dynamo:",set"`
}

type Activity struct {
	ActivityID ActivityID
	Title      string
	Date       string
	Start      string
	End        string
	Confirmed  []UserID `dynamo:",set"`
}

type Availability struct {
	AvailabilityID AvailabilityID
	UserID         UserID
	Date           string
	Start          string
	End            string
}

type Schedule struct {
	ScheduleID int `dynamo:",hash"`
	Year       int
	Month      int
	Day        int
	Tasks      map[string]Task `dynamo:",set"`
}

type Task struct {
	TaskID    int `dynamo:",hash"`
	TaskName  string
	StartTime string //HHMM format string
	EndTime   string //HHMM format string
}
