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
			expected: Log{Meta: Meta{Date: testDate}, Data: []string{"New log entry with tags", "tags: tag0, tag1"}},
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
			expected: `
---
date: Saturday, May 14, 2022
---

22:34	New log entry
`,
		},
		{
			log: Log{Meta: Meta{Date: testDate, Page: "test"}, Data: []string{
				"New log entry",
				"Content",
				"Multiple lines."},
			},
			expected: `
---
date: Saturday, May 14, 2022

page: test
---

22:34	New log entry
Content
Multiple lines.
`,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := formatLog(tt.log)
			if tt.expected != actual {
				t.Fatalf("expected log format: %s\ndid not match actual log format: %s", tt.expected, actual)
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
			log:      Log{Meta: Meta{Date: testDate}},
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
