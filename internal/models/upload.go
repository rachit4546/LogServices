package models

import "context"

type Upload struct {
	ID        string
	Metrics   *Metrics
	Cancel    context.CancelFunc
	Done      chan struct{}
	WorkerCnt int
}
