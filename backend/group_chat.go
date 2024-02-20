package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	chatMessageMinLen = 1
	chatMessageMaxLen = 500
)

// Chat sent over JSON.
type GetChatResponse struct {
	Messages []GetChatResponseMessage `json:"messages"`
	Continue bool                     `json:"continue"`
}

// Message sent over JSON.
type GetChatResponseMessage struct {
	Sender    UserID `json:"sender"`
	Timestamp uint64 `json:"timestamp"`
	Content   string `json:"content"`
}

// New chat sent over JSON.
type PatchChatRequest struct {
	Content string `json:"content"`
}

// API's related to chat within a group.
func RestGroupChatAPI(router *mux.Router, database Database, notification Notification) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if !group.IsMember(user.UserID) {
			http.Error(w, "not a member of group", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			startTimeString := r.URL.Query().Get("start")
			startTime, err := strconv.ParseUint(startTimeString, 10, 64)
			if err != nil {
				startTime = 0
			}
			endTimeString := r.URL.Query().Get("end")
			endTime, err := strconv.ParseUint(endTimeString, 10, 64)
			if err != nil {
				endTime = math.MaxUint64
			}

			messages, more, err := database.ReadMessages(group.GroupID, startTime, endTime)
			if err != nil {
				http.Error(w, "could not read chat", http.StatusInternalServerError)
				return
			}

			var chat GetChatResponse

			chat.Messages = []GetChatResponseMessage{}
			for _, message := range messages {
				chat.Messages = append(chat.Messages, GetChatResponseMessage{
					Sender:    message.Sender,
					Timestamp: message.Timestamp,
					Content:   message.Content,
				})
			}
			chat.Continue = more

			json.NewEncoder(w).Encode(chat)
		case http.MethodPatch:
			var request PatchChatRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if invalidString(w, request.Content, chatMessageMinLen, chatMessageMaxLen) {
				return
			}

			message := Message{
				GroupID:   group.GroupID,
				Sender:    user.UserID,
				Timestamp: unixMillis(),
				Content:   request.Content,
			}
			if err := database.CreateMessage(message); err != nil {
				http.Error(w, fmt.Sprintf("could not create message: %v", err), http.StatusInternalServerError)
				return
			}

			notifyGroup(group, MessageReceived{Message: MessageReceivedMessage{
				GroupID:   message.GroupID,
				Timestamp: message.Timestamp,
				Sender:    message.Sender,
				Content:   message.Content,
			}}, database, notification)
			pushGroup(group, MessagePushed{
				Message: MessagePushedMessage{
					Group:     group.Name,
					Timestamp: message.Timestamp,
					Sender:    user.Name,
					Content:   message.Content,
				},
			}, database)

			WriteJSON(w, nil)
		}
	})
}
