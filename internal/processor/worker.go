package processor

import (
	"context"
	"log-service/internal/models"
)

func worker(ctx context.Context, m *models.Metrics, jobs <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-jobs:
			if !ok {
				return
			}
			parsed, err := ParseLine(line)
			if err != nil || parsed == nil {
				continue
			}
			m.Mu.Lock()
			m.TotalLines++
			m.LevelCount[parsed.Level]++
			m.HourCount[parsed.Hour]++
			m.MsgCount[parsed.Msg]++
			m.Mu.Unlock()
		}
	}
}
