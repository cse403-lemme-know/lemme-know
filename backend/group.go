package main

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// Group sent over JSON.
type GetGroupResponse struct {
	GroupID        GroupID                        `json:"groupId"`
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
	Date       string   `json:"date"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Confirmed  []UserID `json:"confirmed"`
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
		user := Authenticate(w, r, database)
		if user == nil {
			return
		}
		group := Group{
			GroupID: rand.Uint64(),
		}
		if err := database.CreateGroup(group); err != nil {
			http.Error(w, "could not create group", http.StatusInternalServerError)
			return
		}
		WriteJSON(w, group)
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
	RestGroupChatAPI(AddHandler(router, "/chat"), database)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context().Value(UserKey).(*User)
		switch r.Method {
		case http.MethodGet:
			group := r.Context().Value(GroupKey).(*Group)
			WriteJSON(w, GetGroupResponse{
				GroupID: group.GroupID,
				Name:    group.Name,
				Members: group.Members,
				// TODO
				Poll:           nil,
				Availabilities: []GetGroupResponseAvailability{},
				Activities:     []GetGroupResponseActivity{},
			})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
