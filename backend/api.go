package main

import (
	"io"
	"log"
	"net/http"
)

// HTTP multiplexer for the root path.
func Root(database Database, _notification Notification) *http.ServeMux {
	mux := http.NewServeMux()

	// This path won't be reachable via Cloudfront.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Must use GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "\"Hello world!\"")
	})

	AddMux(mux, "/api", Api(database))

	return mux
}

// HTTP multiplexer for the API.
func Api(database Database) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", SessionApi(database))
	return mux
}

// Adds a nested multiplexer at a relative path prefix.
func AddMux(mux *http.ServeMux, prefix string, subMux *http.ServeMux) {
	mux.Handle(prefix+"/", http.StripPrefix(prefix, subMux))
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
