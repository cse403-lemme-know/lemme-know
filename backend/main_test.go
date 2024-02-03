package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// Integration test of HTTP service.
func TestHTTPService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := uint16(30000 + rand.Intn(30000))
	go runLocalService(port, ctx)
	time.Sleep(time.Second / 10)

	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)
	c := &http.Client{
		Jar: jar,
	}

	// Test: create user and get auth cookie.
	response, err := c.Get(fmt.Sprintf("http://localhost:%d/api/user/", port))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(response.Cookies()))
	var getUserResponse GetUserResponse
	MustDecode(t, response.Body, &getUserResponse)
	log.Printf("got user id %d", getUserResponse.UserID)

	// Test: get user again (no new cookie).
	response, err = c.Get(fmt.Sprintf("http://localhost:%d/api/user/", port))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(response.Cookies()))
	var getUserResponse2 GetUserResponse
	MustDecode(t, response.Body, &getUserResponse2)
	assert.Equal(t, getUserResponse.UserID, getUserResponse2.UserID)

	// Test: create group.
	patchGroupRequest := PatchGroupRequest{
		Name: "test",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/", port), patchGroupRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	var patchGroupResponse PatchGroupResponse
	MustDecode(t, response.Body, &patchGroupResponse)
	log.Printf("got group id %d", patchGroupResponse.GroupID)
}

// Integration test of Lambda handler.
func TestLambdaHandler(t *testing.T) {
	type Case struct {
		request  json.RawMessage
		response json.RawMessage
		err      error
	}

	tests := []Case{
		{
			request:  MustMarshal(t, &events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet, Path: "/"}),
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

// Send a patch request via the client.
func Patch(c *http.Client, url string, body any) (resp *http.Response, err error) {
	bodyReader := MustMarshal(&testing.T{}, body)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(bodyReader))
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Send a put request via the client.
func Put(c *http.Client, url string, body any) (resp *http.Response, err error) {
	bodyReader := MustMarshal(&testing.T{}, body)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bodyReader))
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Send a delete request via the client.
func Delete(c *http.Client, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// marshal JSON, failing on any error.
func MustMarshal(t *testing.T, v any) json.RawMessage {
	json, err := json.Marshal(v)
	assert.Nil(t, err)
	return json
}

// unmarshal JSON, failing on any error.
func MustDecode(t *testing.T, r io.Reader, v any) {
	err := json.NewDecoder(r).Decode(v)
	assert.Nil(t, err)
}
