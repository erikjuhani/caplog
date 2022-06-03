package core

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateLog(t *testing.T) {
	testDate := time.Date(2022, 5, 14, 22, 34, 0, 0, time.UTC)

	tests := []struct {
		date     time.Time
		data     string
		expected Log
	}{
		{
			data:     "",
			expected: Log{},
		},
		{
			date:     testDate,
			data:     "New log entry",
			expected: Log{Date: testDate, Data: []string{"New log entry"}},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if !reflect.DeepEqual(tt.expected, CreateLog(tt.date, tt.data)) {
				t.Fatal("expected log did not deeply equal to actual log")
			}
		})
	}
}

func TestFormatLog(t *testing.T) {
	testDate := time.Date(2022, 5, 14, 22, 34, 0, 0, time.UTC)

	tests := []struct {
		log      Log
		expected string
	}{
		{
			log:      Log{Date: testDate, Data: []string{"New log entry"}},
			expected: "22:34\tNew log entry\n",
		},
		{
			log: Log{Date: testDate, Data: []string{
				"New log entry",
				"Content",
				"Multiple lines."},
			},
			expected: "22:34\tNew log entry\nContent\nMultiple lines.\n",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.expected != formatLog(tt.log) {
				t.Fatal("expected log format did not match actual log format")
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	testDate := time.Date(2022, 5, 14, 22, 34, 16, 0, time.UTC)

	tests := []struct {
		log      Log
		expected string
	}{
		{
			log:      Log{},
			expected: "0001-01-01T00:00:00_95c7e5c.log",
		},
		{
			log:      Log{Date: testDate},
			expected: "2022-05-14T22:34:16_1ee35ca.log",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.expected != generateFilename(tt.log) {
				t.Fatal("expected filename did not match actual filename generated")
			}
		})
	}
}
