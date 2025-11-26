package processor

import (
	"context"
	"log-service/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerProcessesLines(t *testing.T) {
	// Create metrics object
	m := &models.Metrics{
		TotalLines: 0,
		LevelCount: map[string]int{},
		HourCount:  map[string]int{},
		MsgCount:   map[string]int{},
	}

	// Prepare jobs channel
	jobs := make(chan string, 10)
	jobs <- "2025-11-25T10:01:23Z INFO User login successful user_id=123"
	jobs <- "2025-11-25T10:02:11Z WARN Rate limit approaching"
	jobs <- "2025-11-25T10:03:55Z ERROR Database timeout query=SELECT * FROM orders"
	close(jobs)

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker
	go worker(ctx, m, jobs)

	// Wait a little for worker to process all lines
	time.Sleep(100 * time.Millisecond)

	// Check metrics
	assert.Equal(t, 3, m.TotalLines)
	assert.Equal(t, 1, m.LevelCount["INFO"])
	assert.Equal(t, 1, m.LevelCount["WARN"])
	assert.Equal(t, 1, m.LevelCount["ERROR"])
	assert.Equal(t, 1, m.MsgCount["User login successful user_id=123"])
	assert.Equal(t, 1, m.MsgCount["Rate limit approaching"])
	assert.Equal(t, 1, m.MsgCount["Database timeout query=SELECT * FROM orders"])
}
