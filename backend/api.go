package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	userMaxConnections = 8
)

// HTTP multiplexer for the root path.
func RestRoot(router *mux.Router, database Database, notification Notification, scheduler Scheduler) {
	RestApi(AddHandler(router, "/api"), database, notification, scheduler)
}

// HTTP multiplexer for the API.
func RestApi(router *mux.Router, database Database, notification Notification, scheduler Scheduler) {
	RestUserAPI(AddHandler(router, "/user"), database, notification)
	RestGroupAPI(AddHandler(router, "/group"), database, notification)
	RestPushAPI(AddHandler(router, "/push"), database, notification)

	// temporary API for testing the scheduler.
	router.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%v", scheduler.Schedule(time.Now().Add(time.Second), nil))))
	})
}

// Adds a nested multiplexer at a relative path prefix.
func AddHandler(router *mux.Router, prefix string) *mux.Router {
	return router.PathPrefix(prefix).Subrouter()
}

// Parse a required path parameter of type uint64.
//
// The second return value is a status flag, set to `true` if ok. If `false`, an error response was sent and the request is done.
func ParseUint64PathParameter(w http.ResponseWriter, r *http.Request, parameterName string) (uint64, bool) {
	parameterString, ok := mux.Vars(r)[parameterName]
	if !ok {
		http.Error(w, "missing "+parameterName, http.StatusBadRequest)
		return 0, false
	}
	parameter, err := strconv.ParseUint(parameterString, 10, 64)
	if err != nil {
		http.Error(w, "invalid "+parameterName, http.StatusBadRequest)
		return 0, false
	}
	return parameter, true
}

// Write HTTP response consisting of JSON.
func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// WebSocket event (connect or disconnect) handler.
//
// UserID is present for connect, and nil for disconnect.
func WebSocket(database Database, connectionID ConnectionID, userID *UserID) error {
	log.Printf("websocket %s userID=%v\n", connectionID, userID)
	if userID == nil {
		userID, err := database.ReadConnection(connectionID)
		if err != nil {
			return err
		}
		if userID == nil {
			return fmt.Errorf("unknown connection")
		}
		if err := database.UpdateUser(*userID, func(user *User) error {
			user.Connections = slices.DeleteFunc(user.Connections, func(c ConnectionID) bool {
				return c == connectionID
			})
			return nil
		}); err != nil {
			return err
		}
		return database.DeleteConnection(connectionID)
	} else {
		_ = database.WriteConnection(connectionID, *userID)
		return database.UpdateUser(*userID, func(user *User) error {
			if len(user.Connections) >= userMaxConnections {
				user.Connections = slices.Delete(user.Connections, 0, 1)
			}
			user.Connections = append(user.Connections, connectionID)
			return nil
		})
	}

}

// Cron event handler.
func Cron(data json.RawMessage) error {
	log.Println("running cron job")
	return nil
}
