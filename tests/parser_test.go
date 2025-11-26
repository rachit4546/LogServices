package tests

import (
	"log-service/internal/processor"
	"testing"
)

func TestParseLine(t *testing.T) {
	line := "INFO 2025-01-01T12:00:00Z User logged in"

	parsed, err := processor.ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed == nil {
		t.Fatalf("expected parsed line, got nil")
	}

	if parsed.Level != "INFO" {
		t.Errorf("expected level INFO, got %s", parsed.Level)
	}
}
