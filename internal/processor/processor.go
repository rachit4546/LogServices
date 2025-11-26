package processor

import (
	"bufio"
	"context"
	"log-service/internal/models"
	"mime/multipart"
	"sync"
)

func StartProcessing(ctx context.Context, up *models.Upload, file multipart.File) {
	jobs := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 0; i < up.WorkerCnt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, up.Metrics, jobs)
		}()
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			close(up.Done)
			return
		default:
			jobs <- scanner.Text()
		}
	}
	close(jobs)
	wg.Wait()
	close(up.Done)
}
