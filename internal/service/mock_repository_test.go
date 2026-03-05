package service

import (
	"expense-tracker/internal/domain"
	"time"
)

type MockRepository struct {
	CreateFunc         func(t domain.Transaction) error
	GetAllFunc         func() ([]domain.Transaction, error)
	GetByDateRangeFunc func(start, end time.Time) ([]domain.Transaction, error)
}

func (m *MockRepository) Create(t domain.Transaction) error {
	return m.CreateFunc(t)
}

func (m *MockRepository) GetAll() ([]domain.Transaction, error) {
	return m.GetAllFunc()
}

func (m *MockRepository) GetByDateRange(start, end time.Time) ([]domain.Transaction, error) {
	return m.GetByDateRangeFunc(start, end)
}
