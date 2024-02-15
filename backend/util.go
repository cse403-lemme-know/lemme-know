package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"unicode/utf8"
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

// Checks a string for validity.
//
// If returns true, error has been sent and should return.
func invalidString(w http.ResponseWriter, input string, minLen uint, maxLen uint) bool {
	if !utf8.ValidString(input) {
		http.Error(w, "invalid utf-8", http.StatusBadRequest)
		return true
	}
	if len(input) < int(minLen) {
		http.Error(w, "too short", http.StatusBadRequest)
		return true
	}
	if len(input) > int(maxLen) {
		http.Error(w, "too long", http.StatusBadRequest)
		return true
	}
	return false
}

// Parses a date like "2006-01-02", returning nil in case of error.
func parseDate(input string) *time.Time {
	t, err := time.Parse(time.DateOnly, input)
	if err == nil {
		return &t
	} else {
		return nil
	}
}

// Parses a time like "15:04", returning nil in case of error.
func parseTime(input string) *time.Time {
	t, err := time.Parse("15:04", input)
	if err == nil {
		return &t
	} else {
		return nil
	}
}

// Checks if a string is a valid date (like "2006-01-02").
//
// If returns true, error has been sent and should return.
func invalidDate(w http.ResponseWriter, input string) bool {
	if parseDate(input) == nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return true
	}
	return false
}

// Checks if a string is a valid time (like "15:04").
//
// If returns true, error has been sent and should return.
func invalidTime(w http.ResponseWriter, input string) bool {
	if parseTime(input) == nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return true
	}
	return false
}
