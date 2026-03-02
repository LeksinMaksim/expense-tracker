package memory

import (
	"expense-tracker/internal/domain"
	"sync"
	"time"
)

type Repository struct {
	mu   sync.RWMutex
	data map[string]domain.Transaction
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]domain.Transaction),
	}
}

func (r *Repository) Create(t domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[t.ID] = t
	return nil
}

func (r *Repository) GetAll() ([]domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions := make([]domain.Transaction, 0, len(r.data))
	for _, t := range r.data {
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *Repository) GetByDateRange(start, end time.Time) ([]domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions := make([]domain.Transaction, 0, len(r.data))

	for _, t := range r.data {
		if (t.Date.Equal(start) || t.Date.After(start)) && t.Date.Before(end) {
			transactions = append(transactions, t)
		}
	}

	return transactions, nil
}
