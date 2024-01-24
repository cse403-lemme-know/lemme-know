package main

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

func SessionApi() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Must use GET", http.StatusMethodNotAllowed)
			return
		}
		// TODO: check session with database.
		_, err := r.Cookie("session")
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: strconv.FormatUint(rand.Uint64(), 36),
			})
		}
		io.WriteString(w, "OK")
	})
	return mux
}
