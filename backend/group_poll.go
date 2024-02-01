package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// New poll sent over JSON.
type PutPollRequest struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

// New votes sent over JSON.
type PatchPollRequest struct {
	Votes []string `json:"votes"`
}

// API's related to polls within a group.
func RestGroupPollAPI(router *mux.Router, database Database) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var request PutPollRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}
			// TODO: database
			_ = request

			WriteJSON(w, nil)
		case http.MethodPatch:
			var request PatchPollRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}
			// TODO: database
			_ = request

			WriteJSON(w, nil)
		case http.MethodDelete:
			// TODO: database
			WriteJSON(w, nil)
		}
	})
}
