package service

import (
	"errors"
	"expense-tracker/internal/domain"
	"testing"
	"time"
)

func TestCreateTransaction_Validation(t *testing.T) {
	mock := &MockRepository{
		CreateFunc: func(tr domain.Transaction) error {
			t.Fatal("Create should not be called when validation fails")
			return nil
		},
	}
	svc := NewExpenseService(mock)

	tests := []struct {
		name    string
		dto     domain.CreateTransactionDTO
		wantErr string
	}{
		{
			name: "amount is zero",
			dto: domain.CreateTransactionDTO{
				Amount:   0,
				Type:     domain.Expense,
				Category: "misc",
			},
			wantErr: "amount must be greater than 0",
		},
		{
			name: "negative amount",
			dto: domain.CreateTransactionDTO{
				Amount:   -10,
				Type:     domain.Income,
				Category: "misc",
			},
			wantErr: "amount must be greater than 0",
		},
		{
			name: "invalid type",
			dto: domain.CreateTransactionDTO{
				Amount:   10,
				Type:     "transfer",
				Category: "misc",
			},
			wantErr: "invalid transaction type",
		},
		{
			name: "empty category",
			dto: domain.CreateTransactionDTO{
				Amount:   10,
				Type:     domain.Expense,
				Category: "",
			},
			wantErr: "category is required",
		},
		{
			name: "invalid date format",
			dto: domain.CreateTransactionDTO{
				Amount:   10,
				Type:     domain.Expense,
				Category: "misc",
				Date:     "2026-03-05",
			},
			wantErr: "invalid date format, use RFC3339",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateTransaction(tc.dto)
			if err == nil {
				t.Fatalf("excepted error %q, got nil", tc.wantErr)
			}
			if err.Error() != tc.wantErr {
				t.Errorf("error = %q; want %q", err.Error(), tc.wantErr)
			}
		})
	}
}

func TestCreateTransaction_Success(t *testing.T) {
	var saved domain.Transaction

	mock := &MockRepository{
		CreateFunc: func(tr domain.Transaction) error {
			saved = tr
			return nil
		},
	}
	svc := NewExpenseService(mock)

	dto := domain.CreateTransactionDTO{
		Amount:      10,
		Type:        domain.Expense,
		Category:    "misc",
		Description: "misc",
		Date:        "2026-03-05T12:00:00Z",
	}

	result, err := svc.CreateTransaction(dto)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}

	if result.ID == "" {
		t.Error("expected non-empty ID")
	}

	if result.Type != domain.Expense {
		t.Errorf("Type = %q; want %q", result.Type, domain.Expense)
	}
	if result.Amount != 10 {
		t.Errorf("Amount = %d; want 10", result.Amount)
	}
	if result.Category != "misc" {
		t.Errorf("Category = %q; want %q", result.Category, "misc")
	}
	if result.Description != "misc" {
		t.Errorf("Description = %q; want %q", result.Description, "misc")
	}
	if saved.ID != result.ID {
		t.Errorf("saved ID = %q; want %q", saved.ID, result.ID)
	}
}

func TestCreateTransaction_DefaultDate(t *testing.T) {
	mock := &MockRepository{
		CreateFunc: func(tr domain.Transaction) error { return nil },
	}
	svc := NewExpenseService(mock)

	before := time.Now().Add(-time.Second)

	dto := domain.CreateTransactionDTO{
		Amount:   10,
		Type:     domain.Expense,
		Category: "misc",
		Date:     "",
	}

	result, err := svc.CreateTransaction(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after := time.Now().Add(time.Second)

	if result.Date.Before(before) || result.Date.After(after) {
		t.Errorf("Date = %v; expected ~now (between %v and %v)", result.Date, before, after)
	}
}

func TestCreateTransaction_RepositoryError(t *testing.T) {
	mock := &MockRepository{
		CreateFunc: func(tr domain.Transaction) error {
			return errors.New("repository error")
		},
	}
	svc := NewExpenseService(mock)

	dto := domain.CreateTransactionDTO{
		Amount:   10,
		Type:     domain.Expense,
		Category: "misc",
		Date:     "2026-03-05T12:00:00Z",
	}

	_, err := svc.CreateTransaction(dto)
	if err == nil {
		t.Fatalf("excepted error, got nil")
	}

	expectedErr := "repository error"
	if err.Error() != expectedErr {
		t.Errorf("error = %q; want %q", err.Error(), expectedErr)
	}
}

func TestGetAllTransaction(t *testing.T) {
	expected := []domain.Transaction{
		{
			ID:       "1",
			Type:     domain.Income,
			Amount:   10,
			Category: "salary",
		},
		{
			ID:       "2",
			Type:     domain.Expense,
			Amount:   100,
			Category: "food",
		},
	}

	mock := &MockRepository{
		GetAllFunc: func() ([]domain.Transaction, error) {
			return expected, nil
		},
	}
	svc := NewExpenseService(mock)

	result, err := svc.GetAllTransactions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(expected) {
		t.Fatalf("got %d transactions; want %d", len(result), len(expected))
	}
}

func TestGetAllTransaction_RepositoryError(t *testing.T) {
	mock := &MockRepository{
		GetAllFunc: func() ([]domain.Transaction, error) {
			return nil, errors.New("repository error")
		},
	}
	svc := NewExpenseService(mock)

	_, err := svc.GetAllTransactions()
	if err == nil {
		t.Fatalf("excepted error, got nil")
	}

	expectedErr := "repository error"
	if err.Error() != expectedErr {
		t.Errorf("error = %q; want %q", err.Error(), expectedErr)
	}
}

func TestGetStatistics(t *testing.T) {
	tests := []struct {
		name             string
		transactions     []domain.Transaction
		wantIncome       int64
		wantExpense      int64
		wantBalance      int64
		wantCategoryFood int64
	}{
		{
			name:             "empty month",
			transactions:     nil,
			wantIncome:       0,
			wantExpense:      0,
			wantBalance:      0,
			wantCategoryFood: 0,
		},
		{
			name: "only income",
			transactions: []domain.Transaction{
				{
					Type:     domain.Income,
					Amount:   10,
					Category: "misc",
				},
			},
			wantIncome:       10,
			wantExpense:      0,
			wantBalance:      10,
			wantCategoryFood: 0,
		},
		{
			name: "mixed income and expense",
			transactions: []domain.Transaction{
				{
					Type:     domain.Income,
					Amount:   100,
					Category: "misc",
				},
				{
					Type:     domain.Expense,
					Amount:   10,
					Category: "food",
				},
				{
					Type:     domain.Expense,
					Amount:   20,
					Category: "food",
				},
				{
					Type:     domain.Expense,
					Amount:   30,
					Category: "misc",
				},
			},
			wantIncome:       100,
			wantExpense:      60,
			wantBalance:      40,
			wantCategoryFood: 30,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &MockRepository{
				GetByDateRangeFunc: func(start, end time.Time) ([]domain.Transaction, error) {
					return tc.transactions, nil
				},
			}
			svc := NewExpenseService(mock)

			summary, err := svc.GetStatistics(time.Now().Year(), time.Now().Month())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if summary.TotalIncome != tc.wantIncome {
				t.Errorf("TotalIncome = %d; want %d", summary.TotalIncome, tc.wantIncome)
			}
			if summary.TotalExpense != tc.wantExpense {
				t.Errorf("TotalExpense = %d; want %d", summary.TotalExpense, tc.wantExpense)
			}
			if summary.Balance != tc.wantBalance {
				t.Errorf("Balance = %d; want %d", summary.Balance, tc.wantBalance)
			}
			if tc.wantCategoryFood > 0 {
				if summary.ExpenseByCategory["food"] != tc.wantCategoryFood {
					t.Errorf("ExpenseByCategory[food] = %d; want %d", summary.ExpenseByCategory["food"], tc.wantCategoryFood)
				}
			}
		})
	}
}

func TestGetStatistics_RepositoryError(t *testing.T) {
	mock := &MockRepository{
		GetByDateRangeFunc: func(start, end time.Time) ([]domain.Transaction, error) {
			return nil, errors.New("repository error")
		},
	}
	svc := NewExpenseService(mock)

	_, err := svc.GetStatistics(time.Now().Year(), time.Now().Month())
	if err == nil {
		t.Fatalf("excepted error, got nil")
	}

	expectedErr := "repository error"
	if err.Error() != expectedErr {
		t.Errorf("error = %q; want %q", err.Error(), expectedErr)
	}
}
