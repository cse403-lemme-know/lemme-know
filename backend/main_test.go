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

type LocalResponseWriter struct {
	response http.Response
	body []byte
}

func (localResponseWriter *LocalResponseWriter) Header() http.Header {
	return localResponseWriter.response.Header
}

func (localResponseWriter LocalResponseWriter) Write(data []byte) (int, error) {
	if localResponseWriter.response.StatusCode == 0 {
		localResponseWriter.WriteHeader
	}
	localResponseWriter.body = append(localResponseWriter.body, ...data)
	return len(data), nil
}

type LocalTransport struct {
	service http.Handler
}

func (localTransport LocalTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var w LocalResponseWriter
	localTransport.service.ServeHTTP(&w, r)
	return &w.response, nil
}

func TestRestAPI(t *testing.T) {
	service := newLocalHandler()

	client := http.Client{
		Transport: LocalTransport{service},
	}
}

func TestLambdaHandler(t *testing.T) {
	type Case struct {
		request  json.RawMessage
		response json.RawMessage
		err      error
	}

	tests := []Case{
		{
			request:  MustMarshal(&events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet, Path: "/"}),
			response: json.RawMessage("404 page not found\n"),
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
