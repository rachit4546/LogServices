package server

import (
	"log-service/internal/api"
	"log-service/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRouter(store *storage.UploadStore) *gin.Engine {
	r := gin.Default()
	h := api.NewHandler(store)
	r.POST("/upload", h.UploadLogs)
	r.GET("/stats/:id", h.GetStats)
	r.POST("/cancel/:id", h.CancelUpload)

	return r
}
