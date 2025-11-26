package processor_test

import (
	"log-service/internal/processor"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLine_ValidLines(t *testing.T) {
	tests := []struct {
		line      string
		wantLevel string
		wantHour  string
		wantMsg   string
	}{
		{
			line:      "2025-11-25T10:01:23Z INFO User login successful user_id=123",
			wantLevel: "INFO",
			wantHour:  "2025-11-25 10:00",
			wantMsg:   "User login successful user_id=123",
		},
		{
			line:      "2025-11-25T15:45:10Z WARN Rate limit approaching",
			wantLevel: "WARN",
			wantHour:  "2025-11-25 15:00",
			wantMsg:   "Rate limit approaching",
		},
		{
			line:      "2025-11-25T23:59:59Z ERROR Database timeout query=SELECT * FROM orders",
			wantLevel: "ERROR",
			wantHour:  "2025-11-25 23:00",
			wantMsg:   "Database timeout query=SELECT * FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			pl, err := processor.ParseLine(tt.line)
			assert.NoError(t, err)
			assert.NotNil(t, pl)
			assert.Equal(t, tt.wantLevel, pl.Level)
			assert.Equal(t, tt.wantHour, pl.Hour)
			assert.Equal(t, tt.wantMsg, pl.Msg)
		})
	}
}

func TestParseLine_InvalidLines(t *testing.T) {
	tests := []struct {
		line string
	}{
		{line: ""},
		{line: "too short"},
		{line: "2025-11-25T10:01:23ZINFO Missing spaces"},
		{line: "invalid-timestamp INFO message"},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			pl, err := processor.ParseLine(tt.line)
			if tt.line == "invalid-timestamp INFO message" {
				// invalid timestamp should return an error
				assert.Error(t, err)
				assert.Nil(t, pl)
			} else {
				// too short or missing parts should return nil without error
				assert.Error(t, err)
				assert.Nil(t, pl)
			}
		})
	}
}

func TestParseLine_CorrectHourFormat(t *testing.T) {
	line := "2025-11-25T14:37:52Z INFO Sample log message"
	pl, err := processor.ParseLine(line)
	assert.NoError(t, err)
	assert.NotNil(t, pl)

	// Verify hour is truncated correctly
	expectedHour := "2025-11-25 14:00"
	assert.Equal(t, expectedHour, pl.Hour)

	// Verify timestamp parsed correctly
	ts, _ := time.Parse(time.RFC3339, "2025-11-25T14:37:52Z")
	assert.Equal(t, ts.Hour(), 14)
}
