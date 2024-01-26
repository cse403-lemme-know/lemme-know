package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API's related to chat within a group.
func RestGroupChatAPI(router *mux.Router, database Database) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			group := r.Context().Value(GroupKey).(*Group)
			startTimeString := r.URL.Query().Get("start")
			startTime, err := strconv.ParseUint(startTimeString, 10, 64)
			if err != nil {
				startTime = 0
			}

			messages, err := database.ReadGroupChat(group.GroupID, startTime)
			if err != nil {
				http.Error(w, "could not read chat", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(messages)
		case http.MethodPut:
			http.Error(w, "not implemented", http.StatusNotImplemented)
		}
	})
}
