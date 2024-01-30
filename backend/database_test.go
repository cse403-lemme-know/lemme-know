package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}

func TestDeleteUser(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}

func TestReadUser(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}

func TestInsertGroup(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}

func TestDeleteGroup(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}

func TestReadGroup(t *testing.T) {
	//Insert user
	table := NewDynamoDB(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	int UserID = 42

	item := User{
		UserID: UserID,
		Username: "Bob",
		GroupIDs: []*string{},
		WebSocketIDs, []*string{},
		Schedules: map[string]Schedule{}
	}

	err := table.CreateUser(item)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assertEquals(t, 1, len(table.users.primarykeys()))
}