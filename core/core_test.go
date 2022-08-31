package core

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateLog(t *testing.T) {
	testDate := time.Date(2022, 5, 14, 22, 34, 0, 0, time.UTC)

	tests := []struct {
		meta     Meta
		data     string
		tags     []string
		expected Log
	}{
		{
			data:     "",
			expected: Log{},
		},
		{
			meta:     Meta{Date: testDate},
			data:     "New log entry",
			expected: Log{Meta: Meta{Date: testDate}, Data: []string{"New log entry"}},
		},
		{
			meta:     Meta{Date: testDate},
			data:     "New log entry with tags",
			tags:     []string{"tag0", "tag1"},
			expected: Log{Meta: Meta{Date: testDate}, Data: []string{"New log entry with tags", "\ntags: tag0, tag1"}},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := NewLog(tt.meta, tt.data, tt.tags)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Fatalf("expected log %v did not equal to actual log %v", tt.expected, actual)
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
			log:      Log{},
			expected: ``,
		},
		{
			log: Log{Meta: Meta{Date: testDate}, Data: []string{"New log entry"}},
			expected: "22:34	New log entry\n",
		},
		{
			log: Log{Meta: Meta{Date: testDate, Page: "test"}, Data: []string{
				"New log entry",
				"Content",
				"Multiple lines."},
			},
			expected: `22:34	New log entry
	Content
	Multiple lines.
`,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := formatLog(tt.log)
			if tt.expected != actual {
				t.Fatalf("expected log format:\n%s\ndid not match actual log format:\n%s", tt.expected, actual)
			}
		})
	}
}

func TestLogFilename(t *testing.T) {
	testDate := time.Date(2022, 5, 14, 22, 34, 16, 0, time.UTC)

	tests := []struct {
		log      Log
		expected string
	}{
		{
			log:      Log{},
			expected: "01-01-0001.log.md",
		},
		{
			log:      Log{Meta: Meta{Date: testDate}},
			expected: "14-05-2022.log.md",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := logFilename(tt.log)
			if tt.expected != actual {
				t.Fatalf("expected filename \"%s\" did not match actual filename \"%s\"", tt.expected, actual)
			}
		})
	}
}
