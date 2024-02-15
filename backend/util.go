package main

import (
	"encoding/json"
	"os"
	"time"
)

// Returns true if and only if executing in an AWS Lambda function.
func isOnLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// Returns Unix time in milliseconds.
func unixMillis() uint64 {
	return uint64(time.Now().UnixMilli())
}

// marshal JSON, failing on any error.
func mustMarshal(v any) json.RawMessage {
	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json
}
