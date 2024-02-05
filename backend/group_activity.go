package main

import (
	"encoding/json"
	"math/rand"
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
		case http.MethodDelete:
			// TODO: database
			WriteJSON(w, nil)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request PatchActivityRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "could not decode body", http.StatusBadRequest)
			return
		}

		group := r.Context().Value(GroupKey).(*Group)

		if err := database.CreateActivity(group.GroupID, Activity{
			ActivityID: rand.Uint64(),
			Title:      request.Title,
			Date:       request.Date,
			Start:      request.Start,
			End:        request.End,
			Confirmed:  []UserID{},
		}); err != nil {
			http.Error(w, "could not create activity", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
