package main

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isInCI() bool {
	return os.Getenv("GITHUB_ACTIONS") != ""
}

func maybeSkip(t *testing.T) {
	if !isInCI() {
		t.Skip("skipping outside of CI")
	}
}

func TestInsertUser(t *testing.T) {
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
	maybeSkip(t)
	//Insert user
	table := NewDynamoDB(nil)

	var userID UserID = rand.Uint64()

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
