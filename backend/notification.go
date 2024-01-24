package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ConnectionID = string

type Notification interface {
	Notify(connectionID ConnectionID, data any) error
}

type APIGateway struct {
	managementAPI *apigatewaymanagementapi.ApiGatewayManagementApi
}

func NewApiGateway(sess *session.Session) *APIGateway {
	managementAPI := apigatewaymanagementapi.New(sess)
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
