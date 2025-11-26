APP_NAME=log-service
CMD_DIR=./cmd/server

.PHONY: all build run test tidy

all: build

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/