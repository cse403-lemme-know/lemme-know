package main

import (
	"net/http"
)

// API multiplexer.
func Api() *http.ServeMux {
	mux := http.NewServeMux()
	AddMux(mux, "/session", SessionApi())
	return mux
}

// Adds a nested multiplexer at a relative path prefix.
func AddMux(mux *http.ServeMux, prefix string, subMux *http.ServeMux) {
	mux.Handle(prefix, http.StripPrefix(prefix, subMux))
}
