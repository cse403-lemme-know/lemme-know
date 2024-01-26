package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RestSpecificGroupAPI(router *mux.Router, database Database) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context().Value(UserKey)
		w.Header().Add("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
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
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(group)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	})
}

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
