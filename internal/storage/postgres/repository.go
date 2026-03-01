package postgres

import (
	"database/sql"
	"expense-tracker/internal/domain"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(t domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, type, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, t.ID, t.Type, t.Amount, t.Category, t.Description, t.Date)
	return err
}

func (r *Repository) GetAll() ([]domain.Transaction, error) {
	query := `SELECT id, type, amount, category, description, date FROM transactions ORDER BY date DESC`
	return r.fetchTransaction(query)
}

func (r *Repository) GetByDateRange(start, end time.Time) ([]domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category, description, date
		FROM transactions
		WHERE date >= $1 AND date < $2
		ORDER BY date DESC
	`

	return r.fetchTransaction(query, start, end)
}

func (r *Repository) fetchTransaction(query string, args ...any) ([]domain.Transaction, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		err := rows.Scan(
			&t.ID,
			&t.Type,
			&t.Amount,
			&t.Category,
			&t.Description,
			&t.Date,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
