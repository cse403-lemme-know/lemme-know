package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	webpush "github.com/Appboy/webpush-go"
)

type Push interface {
	Push(subscription webpush.Subscription) error
}

type WebPush struct {
	privateKey string
	publicKey  string
}

func NewWebPush() *WebPush {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal("unable to access randomness")
	}
	return &WebPush{
		privateKey, publicKey,
	}
}

func (webPush *WebPush) Push(data any, subscription webpush.Subscription) error {
	json, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal notification to JSON: %s", err)
	}
	_, err = webpush.SendNotification(json, &subscription, &webpush.Options{
		VAPIDPrivateKey: webPush.privateKey,
		VAPIDPublicKey:  webPush.publicKey,
		Urgency:         webpush.UrgencyNormal,
		Subscriber:      os.Getenv("DOMAIN"),
		SubIsURL:        true,
	})
	return err
}
