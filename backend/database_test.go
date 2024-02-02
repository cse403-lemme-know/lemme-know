package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	////assert.Equal(t, 1, len(table.users.primarykeys()))
}

func TestDeleteUser(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 1, len(table.users.primarykeys()))

	err = table.DeleteUser(userID)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 0, len(table.users.primarykeys()))
}

func TestReadUser(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 1, len(table.users.primarykeys()))

	user, err := table.ReadUser(userID)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assert.Equal(t, "Bob", user.Name)
}

func TestInsertGroup(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 1, len(table.users.primarykeys()))
}

func TestDeleteGroup(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 1, len(table.users.primarykeys()))
}

func TestReadGroup(t *testing.T) {
	t.Skip("no DynamoDB connection")
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	var userID UserID = 42

	item := User{
		UserID:      userID,
		Name:        "Bob",
		Groups:      []GroupID{},
		Connections: []ConnectionID{},
		Schedules:   map[string]Schedule{},
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	//assert.Equal(t, 1, len(table.users.primarykeys()))
}
