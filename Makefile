.PHONY: build run test lint clean sqlc migrate-up migrate-down

# Build the server binary
build:
	go build -o bin/rss-aggregator ./cmd/server

# Run the server in development
run:
	go run ./cmd/server

# Run all tests
test:
	go test -v -count=1 ./...

# Run linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

# Remove build artifacts
clean:
	rm -rf bin/

# Generate sqlc code from SQL queries
sqlc:
	sqlc generate

# Run database migrations (install: go install github.com/pressly/goose/v3/cmd/goose@latest)
migrate-up:
	goose -dir sql/schema postgres "$$DB_URL" up

migrate-down:
	goose -dir sql/schema postgres "$$DB_URL" down
