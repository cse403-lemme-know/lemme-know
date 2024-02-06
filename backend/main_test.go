package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/websocket"
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

	userID := getUserResponse.UserID
	log.Printf("got user id %d", userID)

	// Test: get user again (no new cookie).
	response, err = c.Get(fmt.Sprintf("http://localhost:%d/api/user/", port))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(response.Cookies()))
	var getUserResponse2 GetUserResponse
	MustDecode(t, response.Body, &getUserResponse2)
	assert.Equal(t, getUserResponse.UserID, getUserResponse2.UserID)

	// Test: open websocket.
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%d", port), Path: "/ws/"}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			log.Printf("received: %s", message)
		}
	}()

	// Test: create group.
	patchGroupRequest := PatchGroupRequest{
		Name: "test",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/", port), patchGroupRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	var patchGroupResponse PatchGroupResponse
	MustDecode(t, response.Body, &patchGroupResponse)

	groupID := patchGroupResponse.GroupID
	log.Printf("got group id %d", groupID)

	// Test: edit group.
	patchGroupRequest = PatchGroupRequest{
		Name:         "test2",
		CalendarMode: "weekdays",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/%d/", port, groupID), patchGroupRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: create poll.
	putPollRequest := PutPollRequest{
		Title:   "who?",
		Options: []string{"me", "you"},
	}
	response, err = Put(c, fmt.Sprintf("http://localhost:%d/api/group/%d/poll/", port, groupID), putPollRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: vote poll.
	patchPollRequest := PatchPollRequest{
		Votes: []string{"me"},
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/%d/poll/", port, groupID), patchPollRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: send chat.
	patchChatRequest := PatchChatRequest{
		Content: "hello",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/%d/chat/", port, groupID), patchChatRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: read chat.
	response, err = c.Get(fmt.Sprintf("http://localhost:%d/api/group/%d/chat/?start=0&end=%s", port, groupID, strconv.FormatUint(math.MaxUint64, 10)))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	var getChatResponse GetChatResponse
	MustDecode(t, response.Body, &getChatResponse)
	assert.Equal(t, false, getChatResponse.Continue)
	assert.Equal(t, 1, len(getChatResponse.Messages))
	assert.Equal(t, userID, getChatResponse.Messages[0].Sender)
	assert.Equal(t, patchChatRequest.Content, getChatResponse.Messages[0].Content)

	// Test: create activity.
	patchActivityRequest := PatchActivityRequest{
		Title: "hang out",
		Date:  "wednesday",
		Start: "1800",
		End:   "1900",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/%d/activity/", port, groupID), patchActivityRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: create availability.
	patchAvailabilityRequest := PatchAvailabilityRequest{
		Date:  "wednesday",
		Start: "1800",
		End:   "1900",
	}
	response, err = Patch(c, fmt.Sprintf("http://localhost:%d/api/group/%d/availability/", port, groupID), patchAvailabilityRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: read group.
	response, err = c.Get(fmt.Sprintf("http://localhost:%d/api/group/%d/", port, groupID))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	var getGroupResponse GetGroupResponse
	MustDecode(t, response.Body, &getGroupResponse)
	assert.Equal(t, patchGroupRequest.Name, getGroupResponse.Name)
	assert.Equal(t, patchGroupRequest.CalendarMode, getGroupResponse.CalendarMode)
	assert.Equal(t, []UserID{userID}, getGroupResponse.Members)
	assert.NotNil(t, getGroupResponse.Poll)
	assert.Equal(t, putPollRequest.Title, getGroupResponse.Poll.Title)
	assert.Equal(t, len(putPollRequest.Options), len(getGroupResponse.Poll.Options))
	for i, option := range putPollRequest.Options {
		assert.Equal(t, option, getGroupResponse.Poll.Options[i].Name)
	}
	assert.Equal(t, 1, len(getGroupResponse.Activities))
	assert.Equal(t, patchActivityRequest.Title, getGroupResponse.Activities[0].Title)
	assert.Equal(t, 1, len(getGroupResponse.Availabilities))
	assert.Equal(t, patchAvailabilityRequest.Date, getGroupResponse.Availabilities[0].Date)
	assert.Equal(t, userID, getGroupResponse.Availabilities[0].UserID)

	activityID := getGroupResponse.Activities[0].ActivityId
	availabilityID := getGroupResponse.Availabilities[0].AvailabilityID
	log.Printf("got activity id %d", activityID)
	log.Printf("got availability id %d", availabilityID)

	// Test: delete activity.
	response, err = Delete(c, fmt.Sprintf("http://localhost:%d/api/group/%d/activity/%d/", port, groupID, activityID))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: delete availability.
	response, err = Delete(c, fmt.Sprintf("http://localhost:%d/api/group/%d/availability/%d/", port, groupID, availabilityID))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: delete poll.
	response, err = Delete(c, fmt.Sprintf("http://localhost:%d/api/group/%d/poll/", port, groupID))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Test: leave/delete group.
	response, err = Delete(c, fmt.Sprintf("http://localhost:%d/api/group/%d/", port, groupID))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
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

// Helper to debug tests.
func _debugBody(r *http.Response) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(string(body))
}
