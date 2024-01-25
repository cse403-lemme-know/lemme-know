package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func MustMarshal(t any) json.RawMessage {
	json, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return json
}

func TestHandler(t *testing.T) {
	type Case struct {
		request  json.RawMessage
		response json.RawMessage
		err      error
	}

	tests := []Case{
		{
			request:  MustMarshal(&events.APIGatewayProxyRequest{Body: "Hi"}),
			response: json.RawMessage("\"Hello world!\""),
			err:      nil,
		},
	}

	for _, test := range tests {
		database := &MockDatabase{}
		notification := &MockNotification{}
		context := context.Background()
		json, err := json.Marshal(test.request)
		if err != nil {
			panic(err)
		}
		response, err := NewHandler(database, notification)(context, json)
		assert.IsType(t, test.err, err)
		assert.Equal(t, string(test.response), response.Body)
	}
}

type MockDatabase struct {
	users  map[UserID]User
	groups map[GroupID]Group
}

func (mockDatabase *MockDatabase) CreateUser(user User) error {
	mockDatabase.users[user.UserID] = user
	return nil
}

func (mockDatabase *MockDatabase) ReadUser(userID UserID) (*User, error) {
	user, ok := mockDatabase.users[userID]
	if ok {
		return &user, nil
	} else {
		return nil, nil
	}
}

func (mockDatabase *MockDatabase) DeleteUser(userID UserID) error {
	delete(mockDatabase.users, userID)
	return nil
}

func (mockDatabase *MockDatabase) CreateGroup(group Group) error {
	mockDatabase.groups[group.GroupID] = group
	return nil
}

func (mockDatabase *MockDatabase) ReadGroup(groupId GroupID) (*Group, error) {
	group, ok := mockDatabase.groups[groupId]
	if ok {
		return &group, nil
	} else {
		return nil, nil
	}
}

func (mockDatabase *MockDatabase) DeleteGroup(groupID GroupID) error {
	delete(mockDatabase.groups, groupID)
	return nil
}

type MockNotification struct{}

func (mockNotification *MockNotification) Notify(connectionID ConnectionID, data any) error {
	return fmt.Errorf("mock notification")
}
