package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
func RestGroupChatAPI(router *mux.Router, database Database) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			group := r.Context().Value(GroupKey).(*Group)
			startTimeString := r.URL.Query().Get("startTime")
			startTime, err := strconv.ParseUint(startTimeString, 10, 64)
			if err != nil {
				startTime = 0
			}
			endTimeString := r.URL.Query().Get("endTime")
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
			user := r.Context().Value(UserKey).(*User)
			group := r.Context().Value(GroupKey).(*Group)

			var request PatchChatRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}
			if err := database.CreateMessage(Message{
				GroupID:   group.GroupID,
				Sender:    user.UserID,
				Timestamp: unixMillis(),
				Content:   request.Content,
			}); err != nil {
				http.Error(w, "could not create message", http.StatusInternalServerError)
				return
			}

			WriteJSON(w, nil)
		}
	})
}
