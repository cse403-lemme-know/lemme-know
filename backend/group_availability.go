package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// New availability sent over JSON.
type PatchAvailabilityRequest struct {
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// API's related to activities within a group.
func RestGroupAvailabilityAPI(router *mux.Router, database Database) {
	router.HandleFunc("/{availabilityID}/", func(w http.ResponseWriter, r *http.Request) {
		_, ok := ParseUint64PathParameter(w, r, "availabilityID")
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

		var request PatchAvailabilityRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "could not decode body", http.StatusBadRequest)
			return
		}

		user := r.Context().Value(GroupKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if err := database.UpdateGroup(group.GroupID, func(group *Group) error {
			group.Availabilities = append(group.Availabilities, Availability{
				AvailabilityID: rand.Uint64(),
				UserID:         user.UserID,
				Date:           request.Date,
				Start:          request.Start,
				End:            request.End,
			})
			return nil
		}); err != nil {
			http.Error(w, "could not create availability", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
