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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	//"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var db = dynamo.New(sess, &aws.Config{Region: aws.String("us-west-2")})
	_ = initializeNewDatabase(db) //Initializing new tables -- redundant code

	_ = apigatewaymanagementapi.New(sess)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	})

	httpHandler := httpadapter.New(http.DefaultServeMux).ProxyWithContext

	lambda.Start(func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
		var ws events.APIGatewayWebsocketProxyRequest
		if err := json.Unmarshal(event, &ws); err == nil {
			log.Println("Received WebSocket event")
			return events.APIGatewayProxyResponse{}, nil
		}

		var http events.APIGatewayProxyRequest
		if err := json.Unmarshal(event, &http); err == nil {
			log.Println("received HTTP request")
			return httpHandler(context, http)
		}

		return events.APIGatewayProxyResponse{}, fmt.Errorf("received unknown message: %s", event)
	})
}
