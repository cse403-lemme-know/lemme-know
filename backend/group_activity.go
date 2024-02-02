package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// New activity sent over JSON.
type PatchActivityRequest struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// API's related to activities within a group.
func RestGroupActivityAPI(router *mux.Router, database Database) {
	router.HandleFunc("/{activityID}/", func(w http.ResponseWriter, r *http.Request) {
		_, ok := ParseUint64PathParameter(w, r, "activityID")
		if !ok {
			return
		}

		switch r.Method {
		case http.MethodPatch:
			var request PatchActivityRequest
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
