package main

import (
	"log"
	"log-service/internal/server"
	"log-service/internal/storage"
)

func main() {
	store := storage.NewUploadStore()
	r := server.SetupRouter(store)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
