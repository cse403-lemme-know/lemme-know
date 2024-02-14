package main

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// New availability sent over JSON.
type PatchAvailabilityRequest struct {
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// API's related to activities within a group.
func RestGroupAvailabilityAPI(router *mux.Router, database Database, notification Notification) {
	router.HandleFunc("/{availabilityID}/", func(w http.ResponseWriter, r *http.Request) {
		availabilityID, ok := ParseUint64PathParameter(w, r, "availabilityID")
		if !ok {
			return
		}

		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if !group.IsMember(user.UserID) {
			http.Error(w, "not a member of group", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			for _, availability := range group.Availabilities {
				if availability.AvailabilityID != availabilityID {
					continue
				}
				if availability.UserID != user.UserID {
					http.Error(w, "cannot delete availability of other member", http.StatusUnauthorized)
					return
				}
			}
			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				group.Availabilities = slices.DeleteFunc(group.Availabilities, func(availability Availability) bool {
					return availability.AvailabilityID == availabilityID && availability.UserID == user.UserID
				})
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not delete availability", http.StatusInternalServerError)
				return
			}
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

		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if !group.IsMember(user.UserID) {
			http.Error(w, "not a member of group", http.StatusUnauthorized)
			return
		}

		if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
			group.Availabilities = append(group.Availabilities, Availability{
				AvailabilityID: GenerateID(),
				UserID:         user.UserID,
				Date:           request.Date,
				Start:          request.Start,
				End:            request.End,
			})
			return nil
		}, database, notification); err != nil {
			http.Error(w, "could not create availability", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
