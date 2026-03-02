include .env
export

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1