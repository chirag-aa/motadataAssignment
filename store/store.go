package store

import (
	"motadataAssignment/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Store stores search records
type Store interface {
	SaveSearch(query string, aiSummary string, articleIDs []string) (models.StoredSearch, error)
	ListSearches() ([]models.StoredSearch, error)
	Clear()
}

// InMemoryStore simple thread-safe store for demos/tests
type InMemoryStore struct {
	mu     sync.RWMutex
	record []models.StoredSearch
}

// NewInMemoryStore returns a new store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		record: []models.StoredSearch{},
	}
}

// SaveSearch stores a search record
func (s *InMemoryStore) SaveSearch(query string, aiSummary string, articleIDs []string) (models.StoredSearch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec := models.StoredSearch{
		ID:                 uuid.NewString(),
		Query:              query,
		AISummaryAnswer:    aiSummary,
		AIRelevantArticles: articleIDs,
		CreatedAt:          time.Now().UTC(),
	}
	s.record = append(s.record, rec)
	return rec, nil
}

func (s *InMemoryStore) ListSearches() ([]models.StoredSearch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.StoredSearch, len(s.record))
	copy(out, s.record)
	return out, nil
}

func (s *InMemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.record = []models.StoredSearch{}
}
