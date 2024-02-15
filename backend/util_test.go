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

func TestCalendarModeParsing(t *testing.T) {
	start, end, dayOfWeek := parseCalendarMode("sus")
	assert.Nil(t, start)
	assert.Nil(t, end)
	assert.False(t, dayOfWeek)
	start, end, dayOfWeek = parseCalendarMode("dayOfWeek")
	assert.Nil(t, start)
	assert.Nil(t, end)
	assert.True(t, dayOfWeek)
	start, end, dayOfWeek = parseCalendarMode("2024-02-15 to 2024-03-04")
	assert.NotNil(t, start)
	assert.NotNil(t, end)
	assert.False(t, dayOfWeek)
}
