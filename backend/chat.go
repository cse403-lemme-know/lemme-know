package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// API's related to chat within a group.
func RestChatAPI(router *mux.Router, database Database) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.Query().Get("start")
		_ = r.URL.Query().Get("end")

		var messages []string

		json.NewEncoder(w).Encode(messages)
	})
}
