package main

import (
	"time"
)

type Group struct {	
	GroupID     int       		     `dynamo:"ID,hash"`	  // Hash key, a.k.a. partition key
	//Time      time.Time // Range key, a.k.a. sort key

	Name   	    string
	//Count     int                  `dynamo:",omitempty"` // Omits if zero value
	Polls  	    []Poll            	 `dynamo:",set"`
	Users  	    map[string]User		 `dynamo:",set"`
}

type User struct {
	UserID      UserID `dynamo:"ID,hash"` // Hash key, a.k.a. partition key
	Name        string
	Groups      []GroupID           `dynamo:",set"`
	Connections []ConnectionID      `dynamo:",set"`
	Schedules   map[string]Schedule `dynamo:",set"`
}

type Poll struct {
	PollID     int                   `dynamo:"ID,hash"` // Hash key, a.k.a. partition key
	Timestamp  time.Time             `dynamo:",range"`  // Range key, a.k.a. sort key
	PollResult map[string]PollResult `dynamo:",set"`
	DoneFlag   bool
}

type Message struct {
	GroupID 	GroupID  		`dynamo:"ID,hash"` //Hash key
	Timestamp 	time.Time		`dynamo:",range"`
	Content   	string			`dynamo:"Message"`
	UserID	  	int
}

type PollResult struct {
	pollResultID int `dynamo:"ID,hash"` //Hash key
	Option       string
	userIDVoted  []int `dynamo:",set"`
}

type Schedule struct {
	ScheduleID int `dynamo:"ID,hash"`
	Year       int
	Month      int
	Day        int
	Tasks      map[string]Task `dynamo:",set"`
}

type Task struct {
	TaskID    int `dynamo:"ID,hash"`
	TaskName  string
	StartTime string //HHMM format string
	EndTime   string //HHMM format string
}
