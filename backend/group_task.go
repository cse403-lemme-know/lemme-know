package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

const (
	taskTitleMinLen = 1
	taskTitleMaxLen = 50
	groupMaxTasks   = 32
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

		if !slices.ContainsFunc(group.Tasks, func(task Task) bool {
			return task.TaskID == taskID
		}) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodPatch:
			var request PatchTaskRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if invalidString(w, request.Title, 0, taskTitleMaxLen) {
				return
			}
			if request.Assignee != nil && !group.IsMember(*request.Assignee) {
				http.Error(w, "assignee not in group", http.StatusBadRequest)
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				for i := range group.Tasks {
					task := &group.Tasks[i]
					if task.TaskID != taskID {
						continue
					}
					if request.Title != "" {
						task.Title = request.Title
					}
					if request.Assignee != nil {
						if !group.IsMember(*request.Assignee) {
							return fmt.Errorf("assignee is not a member")
						}
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

		if invalidString(w, request.Title, taskTitleMinLen, taskTitleMaxLen) {
			return
		}

		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)

		if !group.IsMember(user.UserID) {
			http.Error(w, "not a member of group", http.StatusUnauthorized)
			return
		}

		if invalidAppend(w, group.Tasks, groupMaxTasks) {
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
