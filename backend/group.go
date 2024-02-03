package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// Group sent over JSON.
type GetGroupResponse struct {
	Name           string                         `json:"name"`
	Members        []UserID                       `json:"members,omitempty"`
	Poll           *GetGroupResponsePoll          `json:"poll"`
	Availabilities []GetGroupResponseAvailability `json:"availabilities"`
	Activities     []GetGroupResponseActivity     `json:"activities"`
	CalendarMode   string                         `json:"calendarMode"`
}

// Poll sent over JSON.
type GetGroupResponsePoll struct {
	Options map[string][]UserID `json:"options"`
}

// Availability sent over JSON.
type GetGroupResponseAvailability struct {
	AvailabilityID uint64 `json:"availabilityId"`
	UserID         UserID `json:"userId"`
	Date           string `json:"date"`
	Start          string `json:"start"`
	End            string `json:"end"`
}

// Activity sent over JSON.
type GetGroupResponseActivity struct {
	ActivityId uint64   `json:"activityId"`
	Title      string   `json:"title"`
	Date       string   `json:"date"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Confirmed  []UserID `json:"confirmed"`
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
func RestGroupAPI(router *mux.Router, database Database) {
	router.Use(AuthenticateMiddleware(database))
	RestSpecificGroupAPI(AddHandler(router, "/{groupID}"), database)
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
			GroupID: rand.Uint64(),
			Name:    request.Name,
			Members: []UserID{user.UserID},
		}
		if err := database.CreateGroup(group); err != nil {
			http.Error(w, "could not create group", http.StatusInternalServerError)
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
func RestSpecificGroupAPI(router *mux.Router, database Database) {
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
	RestGroupActivityAPI(AddHandler(router, "/activity"), database)
	RestGroupAvailabilityAPI(AddHandler(router, "/availability"), database)
	RestGroupChatAPI(AddHandler(router, "/chat"), database)
	RestGroupChatAPI(AddHandler(router, "/poll"), database)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context().Value(UserKey).(*User)
		group := r.Context().Value(GroupKey).(*Group)
		switch r.Method {
		case http.MethodGet:
			WriteJSON(w, GetGroupResponse{
				Name:    group.Name,
				Members: group.Members,
				// TODO
				Poll:           nil,
				Availabilities: []GetGroupResponseAvailability{},
				Activities:     []GetGroupResponseActivity{},
			})
		case http.MethodPatch:
			var request PatchGroupRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}

			database.UpdateGroupName(group.GroupID, request.Name)
			// TODO: calendarMode

			WriteJSON(w, nil)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
