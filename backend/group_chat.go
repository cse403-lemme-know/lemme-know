package main

import (
	"encoding/json"
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

			messages, err := database.ReadGroupChat(group.GroupID, startTime)
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

			json.NewEncoder(w).Encode(chat)
		case http.MethodPatch:
			http.Error(w, "not implemented", http.StatusNotImplemented)
		}
	})
}
