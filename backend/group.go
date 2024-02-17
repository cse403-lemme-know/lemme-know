package main

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// Group sent over JSON.
type GetGroupResponse struct {
	Name           string                         `json:"name"`
	Members        []UserID                       `json:"members,omitempty"`
	Poll           *GetGroupResponsePoll          `json:"poll"`
	Availabilities []GetGroupResponseAvailability `json:"availabilities"`
	Activities     []GetGroupResponseActivity     `json:"activities"`
	Tasks          []GetGroupResponseTask         `json:"tasks"`
	CalendarMode   string                         `json:"calendarMode"`
}

// Poll sent over JSON.
type GetGroupResponsePoll struct {
	Title   string                       `json:"title"`
	Options []GetGroupResponsePollOption `json:"options"`
}

// Poll option sent over JSON.
type GetGroupResponsePollOption struct {
	Name  string   `json:"option"`
	Votes []UserID `json:"votes"`
}

// Availability sent over JSON.
type GetGroupResponseAvailability struct {
	AvailabilityID AvailabilityID `json:"availabilityId"`
	UserID         UserID         `json:"userId"`
	Date           string         `json:"date"`
	Start          string         `json:"start"`
	End            string         `json:"end"`
}

// Activity sent over JSON.
type GetGroupResponseActivity struct {
	ActivityID ActivityID `json:"activityId"`
	Title      string     `json:"title"`
	Date       string     `json:"date"`
	Start      string     `json:"start"`
	End        string     `json:"end"`
	Confirmed  []UserID   `json:"confirmed"`
}

// Task sent over JSON.
type GetGroupResponseTask struct {
	TaskID    TaskID `json:"taskId"`
	Title     string `json:"title"`
	Assignee  UserID `json:"assignee"`
	Completed bool   `json:"completed"`
}

// Group properties sent over JSON, used to create or update group.
type PatchGroupRequest struct {
	Name         string `json:"name"`
	CalendarMode string `json:"calendarMode"`
}

// Group ID sent over JSON.
type PatchGroupResponse struct {
	GroupID GroupID `json:"groupId"`
}

// API's related to groups.
func RestGroupAPI(router *mux.Router, database Database, notification Notification) {
	router.Use(AuthenticateMiddleware(database))
	RestSpecificGroupAPI(AddHandler(router, "/{groupID}"), database, notification)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		user := r.Context().Value(UserKey).(*User)
		var request PatchGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "could not decode body", http.StatusBadRequest)
			return
		}
		group := Group{
			GroupID:      GenerateID(),
			Name:         request.Name,
			Members:      []UserID{user.UserID},
			CalendarMode: request.CalendarMode,
		}
		if err := database.CreateGroup(group); err != nil {
			http.Error(w, "could not create group", http.StatusInternalServerError)
			return
		}
		if err := database.UpdateUser(user.UserID, func(user *User) error {
			user.Groups = append(user.Groups, group.GroupID)
			return nil
		}); err != nil {
			http.Error(w, "could not join new group", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, PatchGroupResponse{
			GroupID: group.GroupID,
		})
	})
}

type GroupKeyType struct{}

// Used for looking up group out of request context.
var GroupKey = GroupKeyType(struct{}{})

// API's related to a specific group.
func RestSpecificGroupAPI(router *mux.Router, database Database, notification Notification) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			groupID, ok := ParseUint64PathParameter(w, r, "groupID")
			if !ok {
				return
			}
			group, err := database.ReadGroup(groupID)
			if err != nil {
				http.Error(w, "could not read group", http.StatusInternalServerError)
				return
			}
			if group == nil {
				http.Error(w, "no such group", http.StatusNotFound)
				return
			}
			rWithContext := r.WithContext(context.WithValue(r.Context(), GroupKey, group))
			next.ServeHTTP(w, rWithContext)
		})
	})
	RestGroupActivityAPI(AddHandler(router, "/activity"), database, notification)
	RestGroupAvailabilityAPI(AddHandler(router, "/availability"), database, notification)
	RestGroupChatAPI(AddHandler(router, "/chat"), database, notification)
	RestGroupPollAPI(AddHandler(router, "/poll"), database, notification)
	RestGroupTaskAPI(AddHandler(router, "/task"), database, notification)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)
		switch r.Method {
		case http.MethodGet:
			found := false
			for _, member := range group.Members {
				if member == user.UserID {
					found = true
					break
				}
			}
			if !found {
				if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
					group.Members = append(group.Members, user.UserID)
					return nil
				}, database, notification); err != nil {
					http.Error(w, "could not join group", http.StatusInternalServerError)
					return
				}
				group.Members = append(group.Members, user.UserID)
			}

			response := GetGroupResponse{
				Name:           group.Name,
				CalendarMode:   group.CalendarMode,
				Members:        group.Members,
				Availabilities: []GetGroupResponseAvailability{},
				Activities:     []GetGroupResponseActivity{},
				Tasks:          []GetGroupResponseTask{},
			}

			if group.Poll != nil {
				response.Poll = &GetGroupResponsePoll{
					Title:   group.Poll.Title,
					Options: []GetGroupResponsePollOption{},
				}
				for _, option := range group.Poll.Options {
					response.Poll.Options = append(response.Poll.Options, GetGroupResponsePollOption{
						Name:  option.Name,
						Votes: option.Votes,
					})
				}
			}

			for _, activity := range group.Activities {
				response.Activities = append(response.Activities, GetGroupResponseActivity{
					ActivityID: activity.ActivityID,
					Title:      activity.Title,
					Date:       activity.Date,
					Start:      activity.Start,
					End:        activity.End,
					Confirmed:  activity.Confirmed,
				})
			}

			for _, availability := range group.Availabilities {
				response.Availabilities = append(response.Availabilities, GetGroupResponseAvailability{
					AvailabilityID: availability.AvailabilityID,
					UserID:         availability.UserID,
					Date:           availability.Date,
					Start:          availability.Start,
					End:            availability.End,
				})
			}

			for _, task := range group.Tasks {
				response.Tasks = append(response.Tasks, GetGroupResponseTask{
					TaskID:    task.TaskID,
					Title:     task.Title,
					Assignee:  task.Assignee,
					Completed: task.Completed,
				})
			}

			WriteJSON(w, response)
		case http.MethodPatch:
			if !group.IsMember(user.UserID) {
				http.Error(w, "not a member of group", http.StatusUnauthorized)
				return
			}

			var request PatchGroupRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				if request.Name != "" {
					group.Name = request.Name
				}
				if request.CalendarMode != "" {
					group.CalendarMode = request.CalendarMode
				}
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not update group", http.StatusInternalServerError)
				return
			}

			WriteJSON(w, nil)
		case http.MethodDelete:
			if !group.IsMember(user.UserID) {
				http.Error(w, "not a member of group", http.StatusUnauthorized)
				return
			}
			if err := updateAndNotifyGroup(group.GroupID, func(group *Group) error {
				slices.DeleteFunc(group.Members, func(member UserID) bool {
					return member == user.UserID
				})
				slices.DeleteFunc(group.Availabilities, func(availability Availability) bool {
					return availability.UserID == user.UserID
				})
				for _, activity := range group.Activities {
					slices.DeleteFunc(activity.Confirmed, func(confirmed UserID) bool {
						return confirmed == user.UserID
					})
				}
				// TODO: suboptimal. ideally tasks could be reverted to
				// no assignee, but that complicates the protocol.
				slices.DeleteFunc(group.Tasks, func(task Task) bool {
					return task.Assignee == user.UserID
				})
				if group.Poll != nil {
					for _, option := range group.Poll.Options {
						slices.DeleteFunc(option.Votes, func(vote UserID) bool {
							return vote == user.UserID
						})
					}
				}
				return nil
			}, database, notification); err != nil {
				http.Error(w, "could not leave group (part 1)", http.StatusInternalServerError)
				return
			}
			if err := database.UpdateUser(user.UserID, func(user *User) error {
				slices.DeleteFunc(user.Groups, func(groupID GroupID) bool { return groupID == group.GroupID })
				return nil
			}); err != nil {
				http.Error(w, "could not leave group (part 2)", http.StatusInternalServerError)
				return
			}

			// TODO: delete group if no members left

			WriteJSON(w, nil)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Helper to check if a user is a member of a group.
func (group *Group) IsMember(userID UserID) bool {
	for _, member := range group.Members {
		if member == userID {
			return true
		}
	}
	return false
}
