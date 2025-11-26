package processor_test

import (
	"bytes"
	"context"
	"log-service/internal/models"
	"log-service/internal/processor"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartProcessing(t *testing.T) {
	// Sample log content
	logContent := `2025-11-25T10:01:23Z INFO User login successful user_id=123
2025-11-25T10:02:11Z WARN Rate limit approaching
2025-11-25T10:03:55Z ERROR Database timeout query=SELECT * FROM orders`

	tmp, _ := os.CreateTemp("", "logfile*.log")
	defer os.Remove(tmp.Name())
	tmp.WriteString(logContent)
	tmp.Seek(0, 0)
	file := tmp

	// Setup Upload struct
	up := &models.Upload{
		Metrics:   models.NewMetrics(),
		WorkerCnt: 3,
		Done:      make(chan struct{}),
	}

	ctx := context.Background()

	// Start processing
	go processor.StartProcessing(ctx, up, file)

	// Wait for processing to finish
	select {
	case <-up.Done:
	case <-time.After(1 * time.Second):
		t.Fatal("processing did not finish in time")
	}

	// Check metrics
	assert.Equal(t, 3, up.Metrics.TotalLines)
	assert.Equal(t, 1, up.Metrics.LevelCount["INFO"])
	assert.Equal(t, 1, up.Metrics.LevelCount["WARN"])
	assert.Equal(t, 1, up.Metrics.LevelCount["ERROR"])
	assert.Equal(t, 1, up.Metrics.MsgCount["User login successful user_id=123"])
	assert.Equal(t, 1, up.Metrics.MsgCount["Rate limit approaching"])
	assert.Equal(t, 1, up.Metrics.MsgCount["Database timeout query=SELECT * FROM orders"])
}

func TestStartProcessingWithCancellation(t *testing.T) {
	// Large log content to simulate mid-processing cancellation
	var logContent bytes.Buffer
	for i := 0; i < 1000; i++ {
		logContent.WriteString(`2025-11-25T10:01:23Z INFO User login successful user_id=123` + "\n")
	}

	// 2️⃣ Write to a temporary file
	tmp, err := os.CreateTemp("", "logfile*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name()) // clean up

	// Write the bytes to the temp file
	_, err = tmp.Write(logContent.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// Seek back to beginning before reading
	_, err = tmp.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	// Pass tmp as multipart.File
	file := tmp

	up := &models.Upload{
		Metrics:   models.NewMetrics(),
		WorkerCnt: 5,
		Done:      make(chan struct{}),
	}

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start processing
	go processor.StartProcessing(ctx, up, file)

	// Cancel after a short delay
	time.AfterFunc(5*time.Microsecond, cancel)

	// Wait for done
	select {
	case <-up.Done:
	case <-time.After(1 * time.Second):
		t.Fatal("processing did not finish after cancellation")
	}

	// Ensure some lines were processed, but not necessarily all
	assert.True(t, up.Metrics.TotalLines > 0)
	assert.True(t, up.Metrics.TotalLines < 1000)
}
