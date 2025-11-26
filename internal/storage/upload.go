package storage

import (
	"log-service/internal/models"
	"sync"
)

type UploadStore struct {
	mu   sync.RWMutex
	data map[string]*models.Upload
}

func NewUploadStore() *UploadStore {
	return &UploadStore{
		data: make(map[string]*models.Upload),
	}
}

func (s *UploadStore) Save(id string, u *models.Upload) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = u
}

func (s *UploadStore) Get(id string) (*models.Upload, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}
