package main

import (
	"net/http"
)

// API multiplexer.
func Api(database Database) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", SessionApi(database))
	return mux
}

// Adds a nested multiplexer at a relative path prefix.
func AddMux(mux *http.ServeMux, prefix string, subMux *http.ServeMux) {
	mux.Handle(prefix+"/", http.StripPrefix(prefix, subMux))
}
