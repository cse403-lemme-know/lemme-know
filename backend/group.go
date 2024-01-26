package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API's related to groups.
func RestGroupAPI(router *mux.Router, database Database) {
	router.Use(AuthenticateMiddleware(database))
	RestSpecificGroupAPI(AddHandler(router, "/{groupID}"), database)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
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
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(group)
	})
}

type GroupKeyType struct{}

// Used for looking up group out of request context.
var GroupKey = GroupKeyType(struct{}{})

// API's related to a specific group.
func RestSpecificGroupAPI(router *mux.Router, database Database) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			groupIDString, ok := mux.Vars(r)["groupID"]
			if !ok {
				http.Error(w, "missing group id", http.StatusBadRequest)
				return
			}
			groupID, err := strconv.ParseUint(groupIDString, 10, 64)
			if err != nil {
				http.Error(w, "invalid group id", http.StatusBadRequest)
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
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context().Value(UserKey).(*User)
		w.Header().Add("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			group := r.Context().Value(GroupKey).(*Group)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(group)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
