run:
	go run ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

migrate-up:
	psql "$$PG_DSN" -f db/migrations/0001_init.sql

build:
	go build -o auth-service ./cmd/main.go

clean:
	rm -f auth-service

deps:
	go mod tidy
	go mod download
