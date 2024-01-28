package main

import (
	"encoding/json"
	"log"
	"net/http"

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
