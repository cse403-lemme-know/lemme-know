package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	webpush "github.com/Appboy/webpush-go"
	"github.com/gorilla/mux"
)

/*
func NewWebPush() *WebPush {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal("unable to access randomness")
	}
	return &WebPush{
		privateKey, publicKey,
	}
}
*/

var VAPIDPublicKey string
var VAPIDPrivateKey string

func init() {
	VAPIDPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	VAPIDPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
}

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
			WriteJSON(w, GetPushResponse{
				VAPIDPublicKey: VAPIDPublicKey,
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

func WebPush(data any, subscription webpush.Subscription) error {
	json, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal push to JSON: %w", err)
	}
	resp, err := webpush.SendNotification(json, &subscription, &webpush.Options{
		VAPIDPrivateKey: VAPIDPrivateKey,
		VAPIDPublicKey:  VAPIDPublicKey,
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
