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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/guregu/dynamo"
)

// Take an authenticated AWS session and create a handler for HTTP
// and WebSocket events from AWS API Gateway.
func NewHandler(sess *session.Session) func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	var db = dynamo.New(sess, &aws.Config{Region: aws.String("us-west-2")})
	_ = initializeNewDatabase(db)
	_ = apigatewaymanagementapi.New(sess)

	mux := http.NewServeMux()

	// This path won't be reachable via Cloudfront.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "\"Hello world!\"")
	})

	AddMux(mux, "/api", Api())

	httpHandler := httpadapter.New(mux).ProxyWithContext

	return func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
		var ws events.APIGatewayWebsocketProxyRequest
		if err := json.Unmarshal(event, &ws); err == nil && ws.RequestContext.ConnectionID != "" {
			log.Println("received WebSocket event")
			return events.APIGatewayProxyResponse{}, nil
		}

		var http events.APIGatewayProxyRequest
		if err := json.Unmarshal(event, &http); err == nil {
			log.Println("received HTTP request")
			return httpHandler(context, http)
		}

		return events.APIGatewayProxyResponse{}, fmt.Errorf("received unknown message: %s", event)
	}
}

func main() {
	// Create session from AWS lambda environment.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Start handling events.
	lambda.Start(NewHandler(sess))
}
