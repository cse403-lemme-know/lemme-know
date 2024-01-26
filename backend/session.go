package main

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

func SessionApi(database Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Must use GET", http.StatusMethodNotAllowed)
			return
		}
		user, err := CheckCookie(r, database)
		if err != nil {
			http.Error(w, "could not check cookie", http.StatusInternalServerError)
			return
		}
		if user == nil {
			user := User{
				UserID: rand.Uint64(),
			}
			if err := database.CreateUser(user); err != nil {
				http.Error(w, "could not create user", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "userID",
				Value:    strconv.FormatUint(user.UserID, 10),
				MaxAge:   365 * 24 * 3600,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
			})
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, "\"Ok\"")
	}
}

func CheckCookie(r *http.Request, database Database) (*User, error) {
	cookie, err := r.Cookie("userID")
	if err != nil {
		return nil, nil
	}
	userID, err := strconv.ParseUint(cookie.Value, 10, 64)
	if err != nil {
		return nil, nil
	}
	user, err := database.ReadUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
