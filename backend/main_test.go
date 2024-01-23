package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
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

	// Need deprecated API to avoid missing credential error.
	sess := session.New()
	context := context.Background()
	for _, test := range tests {
		json, err := json.Marshal(test.request)
		if err != nil {
			panic(err)
		}
		response, err := NewHandler(sess)(context, json)
		assert.IsType(t, test.err, err)
		assert.Equal(t, string(test.response), response.Body)
	}
}
