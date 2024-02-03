package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
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
