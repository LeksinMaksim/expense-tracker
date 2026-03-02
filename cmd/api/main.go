package main

import (
	"database/sql"
	"expense-tracker/internal/handlers"
	"expense-tracker/internal/service"
	"expense-tracker/internal/storage/postgres"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/expense_tracker?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}
	defer db.Close()

	err = pingDB(db)
	if err != nil {
		log.Fatalf("Database is not responding: %v", err)
	}
	log.Printf("Successfully connected to PostgreSQL!")

	repo := postgres.NewRepository(db)
	svc := service.NewExpenseService(repo)
	h := handlers.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/transactions", h.CreateTransaction)
	mux.HandleFunc("GET /api/transactions", h.GetTransactions)
	mux.HandleFunc("GET /api/statistics", h.GetStatistics)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func pingDB(db *sql.DB) error {
	var err error
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Database ping failed, retrying in 2 seconds... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}
	return err
}
