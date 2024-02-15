package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	webpush "github.com/Appboy/webpush-go"
	"github.com/gorilla/mux"
)

type GetPushResponse struct {
	VAPIDPublicKey string `json:"vapidPublicKey"`
}

type PatchPushRequest = webpush.Subscription

func RestPushAPI(router *mux.Router, database Database, notification Notification) {
	router.Use(AuthenticateMiddleware(database))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		switch r.Method {
		case http.MethodGet:
			_, vapidPublicKey, err := getVAPIDKeys(database)
			if err != nil {
				http.Error(w, "could not get push keys", http.StatusInternalServerError)
				return
			}
			WriteJSON(w, GetPushResponse{
				VAPIDPublicKey: vapidPublicKey,
			})
		case http.MethodPatch:
			var request PatchPushRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "could not decode body", http.StatusBadRequest)
				return
			}
			user := r.Context().Value(UserKey).(*User)
			if err := database.UpdateUser(user.UserID, func(user *User) error {
				user.Subscriptions = append(user.Subscriptions, request)
				return nil
			}); err != nil {
				http.Error(w, "could not add subscription", http.StatusInternalServerError)
				return
			}
			WriteJSON(w, nil)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func WebPush(data any, subscription webpush.Subscription, database Database) error {
	json, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal push to JSON: %w", err)
	}
	vapidPrivateKey, vapidPublicKey, err := getVAPIDKeys(database)
	if err != nil {
		return fmt.Errorf("could not get VAPID keys: %w", err)
	}
	resp, err := webpush.SendNotification(json, &subscription, &webpush.Options{
		VAPIDPrivateKey: vapidPrivateKey,
		VAPIDPublicKey:  vapidPublicKey,
		Urgency:         webpush.UrgencyNormal,
		Subscriber:      os.Getenv("DOMAIN"),
		SubIsURL:        true,
	})
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}

// For "Voluntary Application server Identification"
//
// Returns private and public VAPID keys, generating and
// storing new ones if needed.
func getVAPIDKeys(database Database) (string, string, error) {
	// Use a single variable so there is never inconsistency
	// between the public and private keys.
	variableName := "VAPID_KEYS"

	keys, err := database.ReadVariable(variableName)
	if err != nil {
		return "", "", err
	}
	if keys == "" || strings.Count(keys, ":") != 1 {
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			return "", "", err
		}
		keys = privateKey + ":" + publicKey
		err = database.WriteVariable(variableName, keys)
		if err != nil {
			// Key is not useful unless we have saved it for next time.
			return "", "", err
		}
	}
	splits := strings.Split(keys, ":")
	// Never out of bounds since we either created the string or
	// ensured it had exactly one delimeter.
	return splits[0], splits[1], nil
}
