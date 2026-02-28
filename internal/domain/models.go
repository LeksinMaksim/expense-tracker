package domain

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          string          `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      int64           `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description,omitempty"`
	Date        time.Time       `json:"date"`
}

type Summary struct {
	TotalIncome       int64            `json:"total_income"`
	TotalExpense      int64            `json:"total_expense"`
	Balance           int64            `json:"balance"`
	ExpenseByCategory map[string]int64 `json:"expense_by_category"`
}

type CreateTransactionDTO struct {
	Type        TransactionType `json:"type"`
	Amount      int64           `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description,omitempty"`
	Date        string          `json:"date"`
}

type TransactionRepository interface {
	Create(t Transaction) error
	GetAll() ([]Transaction, error)
}
