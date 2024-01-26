package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// Take an authenticated AWS session and create a handler for HTTP
// and WebSocket events from AWS API Gateway.
func NewHandler(database Database, _notification Notification) func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {

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

	httpHandler := httpadapter.New(mux).ProxyWithContext

	return func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
		var http events.APIGatewayProxyRequest
		if err := json.Unmarshal(event, &http); err == nil && http.HTTPMethod != "" {
			log.Printf("received HTTP request: %s\n", http.Path)
			return httpHandler(context, http)
		}

		var ws events.APIGatewayWebsocketProxyRequest
		if err := json.Unmarshal(event, &ws); err == nil && ws.RequestContext.ConnectionID != "" {
			log.Printf("received WebSocket event: %s\n", ws.Path)
			return events.APIGatewayProxyResponse{}, nil
		}

		var cron events.EventBridgeEvent
		if err := json.Unmarshal(event, &cron); err == nil && cron.DetailType != "" {
			log.Println("received EventBridge event")
			return events.APIGatewayProxyResponse{}, nil
		}

		return events.APIGatewayProxyResponse{}, fmt.Errorf("received unknown message: %s", event)
	}
}

func main() {
	// Create session from AWS lambda environment.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	database := NewDynamoDB(sess)
	notification := NewApiGateway(sess)

	// Start handling events.
	lambda.Start(NewHandler(database, notification))
}
