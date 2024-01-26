package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/websocket"
)

func NewLambdaHandler(database Database, notification Notification) func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	mux := Root(database, notification)
	httpHandler := httpadapter.New(mux).ProxyWithContext

	return func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
		var http events.APIGatewayProxyRequest
		if err := json.Unmarshal(event, &http); err == nil && http.HTTPMethod != "" {
			log.Printf("received HTTP request: %s\n", http.Path)
			return httpHandler(context, http)
		}

		var ws events.APIGatewayWebsocketProxyRequest
		if err := json.Unmarshal(event, &ws); err == nil && ws.RequestContext.ConnectionID != "" {
			isConnect := ws.RequestContext.EventType == "Connect"
			isDisconnect := ws.RequestContext.EventType == "Disconnect"
			var err error
			if isConnect || isDisconnect {
				err = WebSocket(database, ws.RequestContext.ConnectionID, isConnect)
			}
			return events.APIGatewayProxyResponse{}, err
		}

		var cron events.EventBridgeEvent
		if err := json.Unmarshal(event, &cron); err == nil && cron.DetailType != "" {
			log.Println("received EventBridge event")
			err := Cron()
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{}, fmt.Errorf("received unknown message: %s", event)
	}
}

func RunLambdaService() {
	log.Println("starting AWS lambda service")
	// Create session from AWS lambda environment.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	database := NewDynamoDB(sess)
	notification := NewApiGateway(sess)

	// Start handling events.
	lambda.Start(NewLambdaHandler(database, notification))
}

func RunLocalService() {
	port := 8080
	log.Printf("starting localhost service at http://localhost:%d\n", port)
	database := NewMemoryDatabase()
	notification := NewLocalNotification()

	mux := Root(database, notification)
	upgrader := websocket.Upgrader{} // use default options

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		connectionID := notification.add(c)

		WebSocket(database, connectionID, true)

		go func() {
			defer c.Close()
			defer notification.remove(connectionID)
			defer WebSocket(database, connectionID, false)
			for {
				messageType, _, err := c.ReadMessage()
				if err != nil || messageType == websocket.CloseMessage {
					break
				}
			}
		}()
	})

	// Run `Cron` every hour.
	go func() {
		now := time.Now()
		sleep := 60 - now.Minute()
		time.Sleep(time.Duration(int64(sleep) * int64(time.Minute)))
		Cron()
	}()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}
	log.Fatal(s.ListenAndServe())
}

func IsOnLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

func main() {
	if IsOnLambda() {
		RunLambdaService()
	} else {
		RunLocalService()
	}
}
