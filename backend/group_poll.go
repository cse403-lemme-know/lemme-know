package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

const (
	pollTitleMinLen = 1
	pollTitleMaxLen = 50
	pollMaxOptions  = 4
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
func RestGroupPollAPI(router *mux.Router, database Database, notification Notification) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if !group.IsMember(user.UserID) {
			http.Error(w, "not a member of group", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodPut:
			var request PutPollRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if invalidString(w, request.Title, pollTitleMinLen, pollTitleMaxLen) {
				return
			}

			var options = []PollOption{}

			for _, name := range request.Options {
				if invalidAppend(w, options, pollMaxOptions) {
					return
				}
				options = append(options, PollOption{Name: name, Votes: []UserID{}})
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				group.Poll = &Poll{
					Title:     request.Title,
					Timestamp: unixMillis(),
					Options:   options,
					DoneFlag:  false,
				}
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not create poll", http.StatusInternalServerError)
				return
			}

			WriteJSON(w, nil)
		case http.MethodPatch:
			var request PatchPollRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if group.Poll == nil {
				http.Error(w, "no such poll", http.StatusNotFound)
				return
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				if group.Poll == nil {
					return fmt.Errorf("no such poll")
				}
				for _, option := range group.Poll.Options {
					option.Votes = slices.DeleteFunc(option.Votes, func(o UserID) bool {
						return o == user.UserID
					})
				}
				for _, vote := range request.Votes {
					for _, opt := range group.Poll.Options {
						if opt.Name == vote {
							opt.Votes = append(opt.Votes, user.UserID)
							break
						}
					}
				}
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not vote in poll", http.StatusInternalServerError)
				return
			}

			WriteJSON(w, nil)
		case http.MethodDelete:
			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				group.Poll = nil
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not delete poll", http.StatusInternalServerError)
				return
			}
			WriteJSON(w, nil)
		}
	})
}
