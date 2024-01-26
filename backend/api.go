package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HTTP multiplexer for the root path.
func RestRoot(router *mux.Router, database Database, _notification Notification) {
	RestApi(AddHandler(router, "/api"), database)

	// This path won't be reachable via Cloudfront.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Must use GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "\"Hello world!\"")
	})
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
