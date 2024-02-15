package main

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// New activity sent over JSON.
type PatchActivityRequest struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Confirm *bool  `json:"confirm"`
}

// API's related to activities within a group.
func RestGroupActivityAPI(router *mux.Router, database Database, notification Notification) {
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

		if !slices.ContainsFunc(group.Activities, func(a Activity) bool { return a.ActivityID == activityID }) {
			http.Error(w, "activity not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodPatch:
			var request PatchActivityRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				confirmed := []UserID{}
				if request.Confirm != nil && *request.Confirm {
					confirmed = append(confirmed, user.UserID)
				}
				for i, activity := range group.Activities {
					if activity.ActivityID != activityID {
						continue
					}
					if request.Title != "" && request.Title != activity.Title {
						group.Activities[i].Title = request.Title
					}
					if request.Date != "" && request.Date != activity.Date {
						group.Activities[i].Date = request.Date
						// If the date changes, those that confirmed may no longer be able to attend.
						group.Activities[i].Confirmed = []UserID{}
					}
					if request.Start != "" && request.Start != activity.Start {
						group.Activities[i].Start = request.Start
						// If the start changes, those that confirmed may no longer be able to attend.
						group.Activities[i].Confirmed = []UserID{}
					}
					if request.End != "" && request.End != activity.End {
						group.Activities[i].End = request.End
						// If the end changes, those that confirmed may no longer be able to attend.
						group.Activities[i].Confirmed = []UserID{}
					}
					if request.Confirm != nil {
						if *request.Confirm {
							if !slices.Contains(group.Activities[i].Confirmed, user.UserID) {
								group.Activities[i].Confirmed = append(group.Activities[i].Confirmed, user.UserID)
							}
						} else {
							group.Activities[i].Confirmed = slices.DeleteFunc(group.Activities[i].Confirmed, func(u UserID) bool { return u == user.UserID })
						}
					}
				}

				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not create activity", http.StatusInternalServerError)
				return
			}

			WriteJSON(w, nil)
		case http.MethodDelete:
			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				slices.DeleteFunc(group.Activities, func(activity Activity) bool { return activity.ActivityID == activityID })
				return nil
			}, database, notification); err != nil {
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

		if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
			confirmed := []UserID{}
			if request.Confirm != nil && *request.Confirm {
				confirmed = append(confirmed, user.UserID)
			}
			group.Activities = append(group.Activities, Activity{
				ActivityID: GenerateID(),
				Title:      request.Title,
				Date:       request.Date,
				Start:      request.Start,
				End:        request.End,
				Confirmed:  []UserID{},
			})
			return nil
		}, database, notification); err != nil {
			http.Error(w, "could not create activity", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
