package models

import "sync"

type Metrics struct {
	Mu         sync.RWMutex
	TotalLines int
	LevelCount map[string]int
	HourCount  map[string]int
	MsgCount   map[string]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		LevelCount: make(map[string]int),
		HourCount:  make(map[string]int),
		MsgCount:   make(map[string]int),
	}
}
