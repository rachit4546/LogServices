
# Concurrent Log Processing Service (Go + Gin)
- **Gin** for HTTP API
- **Worker pool** for concurrent log processing


---
## 🚀 Features
- Upload log files (`POST /upload`)
- Process logs concurrently using a worker pool
- Thread-safe metrics aggregation
- Get real-time stats (`GET /stats/{uploadID}`)
- Cancel processing (`POST /cancel/{uploadID}`)

---
## 📁 Project Structure
```
log-service/
├── cmd/server/main.go
├── internal/
│ ├── api/handlers.go
│ ├── processor/
│ │ ├── processor.go
│ │ ├── worker.go
│ │ └── parser.go
│ ├── models/
│ ├── server/router.go
│ └── storage/upload_store.go
├── tests/
│ ├── parser_test.go
│ ├── worker_test.go
│ └── metrics_test.go
├── Makefile
├── go.mod



Running the Service

1. Initialize Module
    go mod tidy

2. Run Server
    make run

Default server runs on :8080

---
API Endpoints


1. POST /upload**
    Upload a log file.

Example using curl:
    curl -X POST -F "file=@logs.txt" http://localhost:8080/upload

Response:
    json : { "uploadID": "d8c1d83a-91ef-4a7b-9c05-4dbd0d4ad311" }

2. GET /stats/{uploadID}
    Retrieve real-time metrics for an upload.
    Response 
    - Total lines processed
    - Count by log level
    - Count by hour
    - Most common messages

3. POST /cancel/{uploadID}**
    Gracefully cancels processing.

Log Format
Every line must follow the format:
    TIMESTAMP LEVEL MESSAGE

Example:
    2025-11-25T10:06:45Z INFO Jobcompleted job_id=abc123 duration=3.4s

Tests
    Run all tests:
    make test
    go test ./...

---
## 📘 Concurrency Model
- Each upload starts a **worker pool** (size configurable per-upload)
- Workers read from a job queue
- A `context.Context` enables cancellation
- Metrics stored per-upload using a `sync.RWMutex`
- No goroutine leakage ensured: all workers exit on cancel or job close