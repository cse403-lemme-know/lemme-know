package db

type group struct {
	GroupID 	int       		`dynamo:"ID,hash"`	  // Hash key, a.k.a. partition key
	//Time   time.Time // Range key, a.k.a. sort key

	GroupName   string
	//Count int                 `dynamo:",omitempty"` // Omits if zero value
	Polls  	    map[string]poll            `dynamo:",set"`
	Users  		map[string]user			 `dynamo:",set"`
	Messages 	map[string]message 		 `dynamo:",set"`
}

type user struct {
	UserID     int       			`dynamo:"ID,hash"`	// Hash key, a.k.a. partition key
	Username   string
	GroupIDs   []string		        `dynamo:",set"`
	WebSocket  string
	Schedules  map[string]schedule	`dynamo:",set"`
}

type poll struct {
	PollID     	int       					`dynamo:"ID,hash"`	// Hash key, a.k.a. partition key
	Timestamp  	time.Time 					`dynamo:",range"`	// Range key, a.k.a. sort key
	PollResult 	map[string]pollResult       `dynamo:",set"`
	DoneFlag   	bool
}

type message struct {
	MessageId 	int  			`dynamo:"ID,hash"` //Hash key
	Timestamp 	time.Time		`dynamo:",range"`
	Content   	string			 `dynamo:"Message"`
	UserID	  	int
}

type pollResult struct {
	pollResultID 	int  		`dynamo:"ID,hash"`	//Hash key
	Option       	string
	userIDVoted	 	[]int		`dynamo:",set"`
}

type schedule struct {
	ScheduleID 		int					`dynamo:"ID,hash"`
	Year 	   		int
	Month	   		int
	Day		   		int
	Tasks	   		map[string]task		`dynamo:",set"`
}

type task struct {
	TaskID			int			`dynamo:"ID,hash"`
	TaskName		string
	StartTime   	string		//HHMM format string
	EndTime			string 		//HHMM format string
}
