package handlers

import (
	"errors"
	"expense-tracker/internal/domain"
	"expense-tracker/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockRepository struct {
	CreateFunc         func(tr domain.Transaction) error
	GetAllFunc         func() ([]domain.Transaction, error)
	GetByDateRangeFunc func(start, end time.Time) ([]domain.Transaction, error)
}

func (m *MockRepository) Create(tr domain.Transaction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(tr)
	}
	return nil
}

func (m *MockRepository) GetAll() ([]domain.Transaction, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockRepository) GetByDateRange(start, end time.Time) ([]domain.Transaction, error) {
	if m.GetByDateRangeFunc != nil {
		return m.GetByDateRangeFunc(start, end)
	}
	return nil, nil
}

func TestHandler_CreateTransaction(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockCreate     func(tr domain.Transaction) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success",
			requestBody: `{"type":"expense","amount":100,"category":"food"}`,
			mockCreate: func(tr domain.Transaction) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"status":"success"`,
		},
		{
			name:           "invalid json",
			requestBody:    `{invalid-json}`,
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `invalid JSON body`,
		},
		{
			name:        "service error",
			requestBody: `{"type":"expense","amount":-10,"category":"food"}`,
			mockCreate: func(tr domain.Transaction) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `amount must be greater than 0`,
		},
		{
			name:        "create repository error",
			requestBody: `{"type":"income","amount":10,"category":"salary"}`,
			mockCreate: func(tr domain.Transaction) error {
				return errors.New("db connection failed")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `db connection failed`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				CreateFunc: tc.mockCreate,
			}
			svc := service.NewExpenseService(mockRepo)
			handler := NewHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CreateTransaction(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("statusCode: %d; want: %d", rr.Code, tc.expectedStatus)
			}
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("body: %q; want: %q", rr.Body.String(), tc.expectedBody)
			}
		})
	}
}

func TestHandler_GetTransactions(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAll     func() ([]domain.Transaction, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success fetching",
			mockGetAll: func() ([]domain.Transaction, error) {
				return []domain.Transaction{
					{
						ID:     "1",
						Type:   domain.Expense,
						Amount: 100,
					},
					{
						ID:     "2",
						Type:   domain.Income,
						Amount: 100,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1"`,
		},
		{
			name: "repository fetch error",
			mockGetAll: func() ([]domain.Transaction, error) {
				return nil, errors.New("db down")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `failed to get transactions`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetAllFunc: tc.mockGetAll,
			}
			svc := service.NewExpenseService(mockRepo)
			handler := NewHandler(svc)

			req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
			rr := httptest.NewRecorder()

			handler.GetTransactions(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("statusCode: %d; want: %d", rr.Code, tc.expectedStatus)
			}
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("body: %q; want: %q", rr.Body.String(), tc.expectedBody)
			}
		})
	}
}

func TestHandler_GetStatistics(t *testing.T) {
	tests := []struct {
		name             string
		urlQuery         string
		mockGetDateRange func(start, end time.Time) ([]domain.Transaction, error)
		expectedStatus   int
		expectedBody     string
	}{
		{
			name:     "success fetching without query",
			urlQuery: "/statistics",
			mockGetDateRange: func(start, end time.Time) ([]domain.Transaction, error) {
				return []domain.Transaction{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"total_income":0`,
		},
		{
			name:     "success fetching with query year and month",
			urlQuery: "/statistics?year=2026&month=3",
			mockGetDateRange: func(start, end time.Time) ([]domain.Transaction, error) {
				if start.Year() != 2026 || start.Month() != time.March {
					return nil, errors.New("wrong dates requested")
				}
				return []domain.Transaction{
					{
						Type:   domain.Income,
						Amount: 1000,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"total_income":1000`,
		},
		{
			name:             "invalid year format",
			urlQuery:         "/statistics?year=not-a-year",
			mockGetDateRange: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     `invalid year format`,
		},
		{
			name:             "invalid month format",
			urlQuery:         "/statistics?year=2026&month=not-a-month",
			mockGetDateRange: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     `invalid month format`,
		},
		{
			name:             "invalid month number",
			urlQuery:         "/statistics?year=2026&month=13",
			mockGetDateRange: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     `invalid month format`,
		},
		{
			name:     "repository error",
			urlQuery: "/statistics",
			mockGetDateRange: func(start, end time.Time) ([]domain.Transaction, error) {
				return nil, errors.New("db offline")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `failed to calculate statistics`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetByDateRangeFunc: tc.mockGetDateRange,
			}
			svc := service.NewExpenseService(mockRepo)
			handler := NewHandler(svc)

			req := httptest.NewRequest(http.MethodGet, tc.urlQuery, nil)
			rr := httptest.NewRecorder()

			handler.GetStatistics(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("statusCode: %d; want: %d", rr.Code, tc.expectedStatus)
			}
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("body: %q; want: %q", rr.Body.String(), tc.expectedBody)
			}

			expectedContentType := "application/json"
			if value := rr.Header().Get("Content-Type"); value != expectedContentType {
				t.Errorf("content-type: %q; want: %q", value, expectedContentType)
			}
		})
	}
}
