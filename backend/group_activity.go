package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"

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
		activityID, ok := ParseUint64PathParameter(w, r, "activityID")
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
			if err := database.UpdateGroup(group.GroupID, func(group *Group) error {
				slices.DeleteFunc(group.Activities, func(activity Activity) bool { return activity.ActivityID == activityID })
				return nil
			}); err != nil {
				http.Error(w, "could not delete activity", http.StatusInternalServerError)
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

		var request PatchActivityRequest
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

		if err := database.UpdateGroup(group.GroupID, func(group *Group) error {
			group.Activities = append(group.Activities, Activity{
				ActivityID: rand.Uint64(),
				Title:      request.Title,
				Date:       request.Date,
				Start:      request.Start,
				End:        request.End,
				Confirmed:  []UserID{},
			})
			return nil
		}); err != nil {
			http.Error(w, "could not create activity", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
