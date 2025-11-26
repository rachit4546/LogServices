package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"log-service/internal/models"
	"log-service/internal/processor"
	"log-service/internal/storage"
	"log-service/utils"
)

type Handler struct {
	Store *storage.UploadStore
}

func NewHandler(store *storage.UploadStore) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) UploadLogs(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}

	uploadID := uuid.New().String()
	ctx, cancel := context.WithCancel(context.Background())
	var val int
	workerCnt := os.Getenv("WORKER_COUNT")
	if workerCnt != "" {
		val, err = strconv.Atoi(workerCnt)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provide Correct WORKER_COUNT"})
		}
	} else {
		fmt.Println("Setting Default value of workerCount to 5")
		val = 5 // default worker count
	}

	upload := &models.Upload{
		ID:        uploadID,
		Metrics:   models.NewMetrics(),
		Cancel:    cancel,
		Done:      make(chan struct{}),
		WorkerCnt: val,
	}
	fmt.Println("Upload = ", &upload)
	h.Store.Save(uploadID, upload)
	go processor.StartProcessing(ctx, upload, src)
	c.JSON(http.StatusOK, gin.H{"uploadID": uploadID})
}

func (h *Handler) GetStats(c *gin.Context) {
	id := c.Param("id")
	upload, ok := h.Store.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	upload.Metrics.Mu.RLock()
	message := utils.SortMap(upload.Metrics.MsgCount)
	resp := gin.H{
		"total":        upload.Metrics.TotalLines,
		"levels":       upload.Metrics.LevelCount,
		"hours":        upload.Metrics.HourCount,
		"top messages": message,
	}
	upload.Metrics.Mu.RUnlock()

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CancelUpload(c *gin.Context) {
	id := c.Param("id")
	upload, ok := h.Store.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	upload.Cancel()
	<-upload.Done

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
