package postgres

import (
	"errors"
	"expense-tracker/internal/domain"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	now := time.Now()

	ts := domain.Transaction{
		ID:          "123",
		Type:        domain.Expense,
		Amount:      10,
		Category:    "misc",
		Description: "misc",
		Date:        now,
	}

	mock.ExpectExec(`INSERT INTO transactions`).WithArgs(ts.ID, ts.Type, ts.Amount, ts.Category, ts.Description, ts.Date).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Create(ts)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "type", "amount", "category", "description", "date"}).AddRow("1", "income", 10, "misc", "misc", now).AddRow("2", "expense", 10, "misc", "misc", now)
	mock.ExpectQuery(`SELECT id, type, amount, category, description, date FROM transactions ORDER BY date DESC`).WillReturnRows(rows)

	transactions, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, transactions, 2)

	assert.Equal(t, "1", transactions[0].ID)
	assert.Equal(t, domain.Income, transactions[0].Type)
	assert.Equal(t, "2", transactions[1].ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRepository_GetByDateRange(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)

	queryRegex := `SELECT id, type, amount, category, description, date FROM transactions WHERE date >= \$1 AND date < \$2 ORDER BY date DESC`
	rows := sqlmock.NewRows([]string{"id", "type", "amount", "category", "description", "date"}).AddRow("3", "expense", 10, "misc", "", start.Add(12*time.Hour))
	mock.ExpectQuery(queryRegex).WithArgs(start, end).WillReturnRows(rows)

	transactions, err := repo.GetByDateRange(start, end)
	require.NoError(t, err)
	require.Len(t, transactions, 1)

	assert.Equal(t, "3", transactions[0].ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRepository_FetchTransaction_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	mock.ExpectQuery(`SELECT`).WillReturnError(errors.New("query error"))

	_, err = repo.GetAll()
	require.Error(t, err)
	assert.Equal(t, "query error", err.Error())
}

func TestRepository_FetchTransaction_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	_, err = repo.GetAll()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected")
}

func TestRepository_FetchTransaction_RowError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id", "type", "amount", "category", "description", "date"}).AddRow("1", "income", 10, "misc", "misc", time.Now()).RowError(0, errors.New("row error"))
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	_, err = repo.GetAll()
	require.Error(t, err)
	assert.Equal(t, "row error", err.Error())
}
