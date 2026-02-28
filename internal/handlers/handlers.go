package handlers

import (
	"encoding/json"
	"expense-tracker/internal/domain"
	"expense-tracker/internal/service"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *service.ExpenseService
}

func NewHandler(service *service.ExpenseService) *Handler {
	return &Handler{
		service: service,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateTransactionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	transaction, err := h.service.CreateTransaction(dto)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":     transaction.ID,
		"status": "success",
	})
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get transactions")
		return
	}

	writeJSON(w, http.StatusOK, transactions)
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		} else {
			writeError(w, http.StatusBadRequest, "invalid year format")
			return
		}
	}

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = time.Month(m)
		} else {
			writeError(w, http.StatusBadRequest, "invalid month format")
			return
		}
	}

	summary, err := h.service.GetStatistics(year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to calculate statistics")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
