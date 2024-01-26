package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RestGroupAPI(router *mux.Router, database Database) {
	AddHandler(router, "/{groupID}").HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("bar")
		user := Authenticate(w, r, database)
		if user == nil {
			return
		}
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
		case http.MethodPut:
			group := Group{
				GroupID: rand.Uint64(),
			}
			if err := database.CreateGroup(group); err != nil {
				http.Error(w, "could not create group", http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(group)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
