package service

import (
	"errors"
	"expense-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ExpenseService struct {
	repo domain.TransactionRepository
}

func NewExpenseService(repo domain.TransactionRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) CreateTransaction(dto domain.CreateTransactionDTO) (domain.Transaction, error) {
	if dto.Amount <= 0 {
		return domain.Transaction{}, errors.New("amount must be greater than 0")
	}
	if dto.Type != domain.Income && dto.Type != domain.Expense {
		return domain.Transaction{}, errors.New("invalid transaction type")
	}
	if dto.Category == "" {
		return domain.Transaction{}, errors.New("category is required")
	}

	var transDate time.Time
	if dto.Date == "" {
		transDate = time.Now()
	} else {
		parsed, err := time.Parse(time.RFC3339, dto.Date)
		if err != nil {
			return domain.Transaction{}, errors.New("invalid date format, use RFC3339")
		}
		transDate = parsed
	}

	transaction := domain.Transaction{
		ID:          uuid.New().String(),
		Type:        dto.Type,
		Amount:      dto.Amount,
		Category:    dto.Category,
		Description: dto.Description,
		Date:        transDate,
	}

	if err := s.repo.Create(transaction); err != nil {
		return domain.Transaction{}, err
	}

	return transaction, nil
}

func (s *ExpenseService) GetAllTransactions() ([]domain.Transaction, error) {
	return s.repo.GetAll()
}

func (s *ExpenseService) GetStatistics(year int, month time.Month) (domain.Summary, error) {
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	transactions, err := s.repo.GetByDateRange(startOfMonth, startOfNextMonth)
	if err != nil {
		return domain.Summary{}, err
	}

	summary := domain.Summary{
		ExpenseByCategory: make(map[string]int64),
	}

	for _, transaction := range transactions {
		if transaction.Type == domain.Income {
			summary.TotalIncome += transaction.Amount
		} else if transaction.Type == domain.Expense {
			summary.TotalExpense += transaction.Amount
			summary.ExpenseByCategory[transaction.Category] += transaction.Amount
		}
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}
