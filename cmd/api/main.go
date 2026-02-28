package main

import (
	"expense-tracker/internal/handlers"
	"expense-tracker/internal/service"
	"expense-tracker/internal/storage/memory"
	"log"
	"net/http"
)

func main() {
	repo := memory.NewRepository()
	svc := service.NewExpenseService(repo)
	h := handlers.NewHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/transactions", h.CreateTransaction)
	mux.HandleFunc("GET /api/transactions", h.GetTransactions)
	mux.HandleFunc("GET /api/statistics", h.GetStatistics)

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
