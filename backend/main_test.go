package main

import (
	"context"
	"encoding/json"
	"net/http"
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
			request:  MustMarshal(&events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet, Path: "/"}),
			response: json.RawMessage("\"Hello world!\"\n"),
			err:      nil,
		},
	}

	for _, test := range tests {
		database := NewMemoryDatabase()
		notification := NewLocalNotification()
		context := context.Background()
		json, err := json.Marshal(test.request)
		if err != nil {
			panic(err)
		}
		response, err := newLambdaHandler(database, notification)(context, json)
		assert.IsType(t, test.err, err)
		assert.Equal(t, string(test.response), response.Body)
	}
}
