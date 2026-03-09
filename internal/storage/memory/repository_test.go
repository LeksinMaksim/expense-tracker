package memory

import (
	"expense-tracker/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateAndGetAll(t *testing.T) {
	repo := NewRepository()

	now := time.Now()
	ts1 := domain.Transaction{
		ID:          "1",
		Type:        domain.Income,
		Amount:      10,
		Category:    "misc",
		Description: "misc",
		Date:        now,
	}
	ts2 := domain.Transaction{
		ID:          "2",
		Type:        domain.Expense,
		Amount:      10,
		Category:    "misc",
		Description: "misc",
		Date:        now.Add(time.Hour),
	}

	err := repo.Create(ts1)
	require.NoError(t, err)

	err = repo.Create(ts2)
	require.NoError(t, err)

	transactions, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, transactions, 2)

	tsMap := make(map[string]domain.Transaction)
	for _, ts := range transactions {
		tsMap[ts.ID] = ts
	}

	assert.Equal(t, ts1, tsMap["1"])
	assert.Equal(t, ts2, tsMap["2"])
}

func TestRepository_GetByDateRange(t *testing.T) {
	repo := NewRepository()
	baseTime := time.Date(2026, time.March, 1, 10, 0, 0, 0, time.UTC)

	err := repo.Create(domain.Transaction{
		ID:   "1",
		Date: baseTime.Add(-24 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.Create(domain.Transaction{
		ID:   "2",
		Date: baseTime.Add(time.Hour),
	})
	require.NoError(t, err)

	err = repo.Create(domain.Transaction{
		ID:   "3",
		Date: baseTime.Add(2 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.Create(domain.Transaction{
		ID:   "4",
		Date: baseTime.Add(24 * time.Hour),
	})
	require.NoError(t, err)

	start := baseTime
	end := baseTime.Add(24 * time.Hour)

	transactions, err := repo.GetByDateRange(start, end)
	require.NoError(t, err)
	assert.Len(t, transactions, 2)

	tsMap := make(map[string]bool)
	for _, ts := range transactions {
		tsMap[ts.ID] = true
	}

	assert.True(t, tsMap["2"])
	assert.True(t, tsMap["3"])
}
