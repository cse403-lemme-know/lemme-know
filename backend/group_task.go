package main

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// New/updated task sent over JSON.
type PatchTaskRequest struct {
	Title     string  `json:"title"`
	Assignee  *UserID `json:"assignee"`
	Completed *bool   `json:"completed"`
}

// API's related to activities within a group.
func RestGroupTaskAPI(router *mux.Router, database Database, notification Notification) {
	router.HandleFunc("/{taskID}/", func(w http.ResponseWriter, r *http.Request) {
		taskID, ok := ParseUint64PathParameter(w, r, "taskID")
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
		case http.MethodPatch:
			var request PatchTaskRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				for _, task := range group.Tasks {
					if task.TaskID != taskID {
						continue
					}
					if request.Title != "" {
						task.Title = request.Title
					}
					if request.Assignee != nil {
						task.Assignee = *request.Assignee
					}
					if request.Completed != nil {
						task.Completed = *request.Completed
					}
				}
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not update task", http.StatusInternalServerError)
				return
			}
			WriteJSON(w, nil)
		case http.MethodDelete:
			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				group.Tasks = slices.DeleteFunc(group.Tasks, func(task Task) bool {
					return task.TaskID == taskID
				})
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not delete task", http.StatusInternalServerError)
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

		var request PatchTaskRequest
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

		completed := false
		assignee := user.UserID
		if request.Completed != nil {
			completed = *request.Completed
		}
		if request.Assignee != nil && group.IsMember(*request.Assignee) {
			assignee = *request.Assignee
		}

		if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
			group.Tasks = append(group.Tasks, Task{TaskID: GenerateID(), Title: request.Title, Completed: completed, Assignee: assignee})
			return nil
		}, database, notification); err != nil {
			http.Error(w, "could not create task", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, nil)
	})
}
