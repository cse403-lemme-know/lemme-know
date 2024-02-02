package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HTTP multiplexer for the root path.
func RestRoot(router *mux.Router, database Database, _notification Notification) {
	RestApi(AddHandler(router, "/api"), database)
}

// HTTP multiplexer for the API.
func RestApi(router *mux.Router, database Database) {
	RestUserAPI(AddHandler(router, "/user"), database)
	RestGroupAPI(AddHandler(router, "/group"), database)
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
func WebSocket(database Database, connectionId ConnectionID, isConnect bool) error {
	log.Printf("websocket %s isConnect=%t\n", connectionId, isConnect)
	return nil
}

// Cron event handler.
func Cron() error {
	log.Println("running cron job")
	return nil
}
