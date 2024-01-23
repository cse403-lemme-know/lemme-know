package db

type group struct {
	GroupID 	int       // Hash key, a.k.a. partition key
	//Time   time.Time // Range key, a.k.a. sort key

	GroupName   string
	//Count int                 `dynamo:",omitempty"` // Omits if zero value
	Polls  	    []poll            `dynamo:",set"`
	Users  		[]user			 `dynamo:",set"`
	Messages 	[]message 		 `dynamo:",set"`
}

type user struct {
	UserID     int       // Hash key, a.k.a. partition key
	//Time   time.Time // Range key, a.k.a. sort key

	Username   string               `dynamo:"Message"`    // Change name in the database
	GroupIDs   []string		        `dynamo:",set"`
	Schedules  []schedule			`dynamo:",set"`
}

type poll struct {
	PollID     	int       // Hash key, a.k.a. partition key
	Timestamp  	time.Time // Range key, a.k.a. sort key
	PollResult 	[]pollResult       `dynamo:",set"`
	DoneFlag   	bool
}

type message struct {
	MessageId 	int  //Hash key
	
	Timestamp 	time.Time
	Content   	string			 `dynamo:"Message"`
	UserID	  	int
}

type pollResult struct {
	pollResultID 	int  //Hash key

	Option       	string
	userIDVoted	 	[]int		`dynamo:",set"`
}

type schedule struct {
	ScheduleID 		int
	Year 	   		int
	Month	   		int
	Day		   		int
	Tasks	   		[]task		`dynamo:",set"`
}

type task struct {
	TaskID			int
	TaskName		string
	StartTime   	string		//HHMM format string
	EndTime			string 		//HHMM format string
}
