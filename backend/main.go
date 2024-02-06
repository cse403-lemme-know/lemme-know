package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Creates a function that can handle JSON events from AWS services.
func newLambdaHandler(database Database, notification Notification) func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	router := mux.NewRouter()

	// Expose the entire Rest API.
	RestRoot(router, database, notification)

	handler := applyCors(router)

	// Convert Go http handler to AWS Lambda http handler.
	httpHandler := httpadapter.New(handler).ProxyWithContext

	return func(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
		// Check if the event is an AWS API Gateway HTTP Rest request.
		var rest events.APIGatewayProxyRequest
		if err := json.Unmarshal(event, &rest); err == nil && rest.HTTPMethod != "" {
			log.Printf("received HTTP request: %s\n", rest.Path)
			return httpHandler(context, rest)
		}

		// Check if the event is an AWS API Gateway HTTP WebSocket event.
		var ws events.APIGatewayWebsocketProxyRequest
		if err := json.Unmarshal(event, &ws); err == nil && ws.RequestContext.ConnectionID != "" {
			// Construct a request to access cookies.
			request, err := http.NewRequest(
				http.MethodGet,
				"http://example.com",
				strings.NewReader(ws.Body),
			)
			if err != nil {
				// Unreachable.
				panic(err)
			}
			for key, value := range ws.Headers {
				request.Header.Add(key, value)
			}
			for key, values := range ws.MultiValueHeaders {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
			user, err := CheckCookie(request, database)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}
			if user == nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusUnauthorized,
					Body:       "no such user",
				}, nil
			}

			isConnect := ws.RequestContext.EventType == "Connect"
			isDisconnect := ws.RequestContext.EventType == "Disconnect"
			if isConnect || isDisconnect {
				err = WebSocket(database, ws.RequestContext.ConnectionID, user.UserID, isConnect)
			}
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       "",
			}, err
		}

		// Check if the event is an AWS EventBridge cron event.
		var cron events.EventBridgeEvent
		if err := json.Unmarshal(event, &cron); err == nil && cron.DetailType != "" {
			log.Println("received EventBridge event")
			err := Cron()
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{}, fmt.Errorf("received unknown message: %s", event)
	}
}

// Handle events forever within AWS Lambda with a non-volatile database.
func runLambdaService() {
	log.Println("starting AWS lambda service")
	// Create session from AWS lambda environment.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	database := NewDynamoDB(sess)
	notification := NewApiGateway(sess)

	// Start handling events forever.
	lambda.Start(newLambdaHandler(database, notification))
}

// Handle events forever on localhost with a volatile database.
//
// Returns errors except if they were due to ctx being canceled.
func runLocalService(port uint16, ctx context.Context) error {
	database := NewMemoryDatabase()
	notification := NewLocalNotification()

	router := mux.NewRouter()
	upgrader := websocket.Upgrader{} // use default options

	// In addition to the Rest API, expose WebSocket capabilities.
	router.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckCookie(r, database)
		if err != nil {
			http.Error(w, "could not check cookie", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "no such user", http.StatusUnauthorized)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not open WebSocket", http.StatusUpgradeRequired)
			return
		}

		connectionID := notification.add(c)

		WebSocket(database, connectionID, user.UserID, true)

		go func() {
			defer c.Close()
			defer notification.remove(connectionID)
			defer WebSocket(database, connectionID, user.UserID, false)
			for {
				messageType, _, err := c.ReadMessage()
				if err != nil || messageType == websocket.CloseMessage {
					break
				}
			}
		}()
	})

	// Expose the Rest API.
	RestRoot(router, database, notification)

	// In addition to the Rest API, reverse proxy to the development client's origin server.
	clientOrigin, err := url.Parse("http://localhost:5173/")
	if err != nil {
		panic(err)
	}
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set req Host, URL and Request URI to forward a request to the origin server
		r.Host = clientOrigin.Host
		r.URL.Host = clientOrigin.Host
		r.URL.Scheme = clientOrigin.Scheme
		r.RequestURI = ""

		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		for name, values := range originServerResponse.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(originServerResponse.StatusCode)
		io.Copy(w, originServerResponse.Body)
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        applyCors(router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 10,
		BaseContext:    func(net.Listener) context.Context { return ctx },
	}

	// Run cron job every hour. If ctx cancelled, shut down the server.
	go func() {
		now := time.Now()
		sleep := 60 - now.Minute()
		select {
		case <-ctx.Done():
			// ctx cancelled
			_ = s.Close()
			return
		case <-time.After(time.Duration(int64(sleep) * int64(time.Minute))):
			Cron()
		}
	}()

	// Serve HTTP until it encounters an error or the context is canceled.
	err = s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Allow exotic HTTP methods, credentials.
func applyCors(handler http.Handler) http.Handler {
	return handlers.CORS(handlers.AllowCredentials(), handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}))(handler)
}

// Returns true if and only if executing in an AWS Lambda function.
func isOnLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// Returns Unix time in milliseconds.
func unixMillis() uint64 {
	return uint64(time.Now().UnixMilli())
}

func main() {
	if isOnLambda() {
		runLambdaService()
	} else {
		const port = 8080
		log.Printf("starting localhost service at http://localhost:%d\n", port)
		runLocalService(port, context.Background())
	}
}
