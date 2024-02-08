package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/gorilla/websocket"
)

type ConnectionID = string

type GroupChanged struct {
	Group GroupChangedGroup `json:"group"`
}

type GroupChangedGroup struct {
	GroupID GroupID `json:"groupID"`
}

// Send a best-effort notification to all group members.
//
// If `data` is `nil`, then just send a group-changed notification.
func notifyGroup(group *Group, data any, database Database, notification Notification) {
	dataOrGroupChanged := data
	if dataOrGroupChanged == nil {
		dataOrGroupChanged = GroupChanged{
			Group: GroupChangedGroup{
				GroupID: group.GroupID,
			},
		}
	}
	var wait sync.WaitGroup
	for _, userID := range group.Members {
		wait.Add(1)
		go func() {
			defer wait.Done()
			user, err := database.ReadUser(userID)
			if err != nil || user == nil {
				// Ignore errors as notification is best-effort.
				return
			}
			for _, connectionID := range user.Connections {
				// Ignore errors as notification is best-effort.
				_ = notification.Notify(connectionID, dataOrGroupChanged)
			}
		}()
	}
	wait.Wait()
}

// Update a group (like `Database.UpdateGroup`) and notify the members
// (like `Notification>NotifyGroup`)
//
// Errors are passed through from `Database.UpdateGroup`. Notification is
// skipped in the case of an error.
func updateAndNotifyGroup(groupID GroupID, transaction func(*Group) error, database Database, notification Notification) error {
	var g *Group = nil
	err := database.UpdateGroup(groupID, func(group *Group) error {
		g = group
		return transaction(group)
	})
	if err != nil {
		return err
	}
	notifyGroup(g, nil, database, notification)
	return nil
}

// A service capable of notifying clients regardless of how many tab(s) they have open.
type Notification interface {
	// Sends a JSON notification to the client with the provided connection ID.
	//
	// Returns a non-nil error if there is evidence to suggest that delivery failed,
	// allthough the result may always be inconclusive.
	Notify(connectionID ConnectionID, data any) error
}

// An AWS service to send a message on an open WebSocket.
type APIGateway struct {
	managementAPI *apigatewaymanagementapi.ApiGatewayManagementApi
}

func NewApiGateway(sess *session.Session) *APIGateway {
	config := aws.NewConfig()
	if endpoint := os.Getenv("AWS_API_GATEWAY_WS_ENDPOINT"); endpoint != "" {
		config.WithEndpoint(endpoint)
	} else {
		log.Println("could not get ws endpoint from environment so ws notifications will not work")
	}
	managementAPI := apigatewaymanagementapi.New(sess, config)
	return &APIGateway{managementAPI}
}

func (apiGateway *APIGateway) Notify(connectionID ConnectionID, data any) error {
	json, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal notification to JSON: %s", err)
	}
	_, err = apiGateway.managementAPI.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionID,
		Data:         json,
	})
	if err != nil {
		return fmt.Errorf("could not notify: %s", err)
	}
	return nil
}

// A local service currently incapable of actually notifying anything.
type LocalNotification struct {
	websockets map[string]*websocket.Conn
	mu         sync.Mutex
}

func NewLocalNotification() *LocalNotification {
	return &LocalNotification{
		websockets: make(map[string]*websocket.Conn),
		mu:         sync.Mutex{},
	}
}

func (localNotification *LocalNotification) Notify(connectionID ConnectionID, data any) error {
	localNotification.mu.Lock()
	defer localNotification.mu.Unlock()
	var err error
	if conn := localNotification.websockets[connectionID]; conn != nil {
		err = conn.WriteJSON(data)
	} else {
		err = fmt.Errorf("no such connection")
	}
	return err
}

func (localNotification *LocalNotification) add(conn *websocket.Conn) ConnectionID {
	localNotification.mu.Lock()
	defer localNotification.mu.Unlock()
	var connectionID string
	for connectionID == "" || localNotification.websockets[connectionID] != nil {
		connectionID = strconv.FormatUint(GenerateID(), 10)
	}
	localNotification.websockets[connectionID] = conn
	return connectionID
}

func (localNotification *LocalNotification) remove(connectionID ConnectionID) {
	localNotification.mu.Lock()
	defer localNotification.mu.Unlock()
	delete(localNotification.websockets, connectionID)
}
