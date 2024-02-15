package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateParsing(t *testing.T) {
	assert.Nil(t, parseDate("sus"))
	assert.NotNil(t, parseDate("2024-02-15"))
}

func TestTimeParsing(t *testing.T) {
	assert.Nil(t, parseTime("sus"))
	assert.NotNil(t, parseTime("19:10"))
}
