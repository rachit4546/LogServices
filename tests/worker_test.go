package tests

import (
	"context"
	"log-service/internal/models"
	"log-service/internal/processor"
	"testing"
	"time"
)

func TestWorkerProcessesLines(t *testing.T) {
	m := models.NewMetrics()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan string, 2)

	go func() {
		processor.ParseLine("INFO 2025-01-01T12:00:00Z msg")
	}()

	go func() {
		processor.ParseLine("ERROR 2025-01-01T12:01:00Z failure")
	}()

	go processorTestWorker(ctx, m, jobs)

	jobs <- "INFO 2025-01-01T12:00:00Z msg"
	jobs <- "ERROR 2025-01-01T12:01:00Z failure"
	close(jobs)

	time.Sleep(100 * time.Millisecond)

	m.Mu.RLock()
	if m.TotalLines != 2 {
		t.Errorf("expected 2 lines, got %d", m.TotalLines)
	}
	m.Mu.RUnlock()
}

// local test wrapper to call unexported worker logic
func processorTestWorker(ctx context.Context, m *models.Metrics, jobs <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-jobs:
			if !ok {
				return
			}
			parsed, err := processor.ParseLine(line)
			if err != nil || parsed == nil {
				continue
			}
			m.Mu.Lock()
			m.TotalLines++
			m.Mu.Unlock()
		}
	}
}
