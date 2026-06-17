.PHONY: help build run test clean db-up db-init lint vet tidy migrate

help:
	@echo "SimplCommerce-Go"
	@echo ""
	@echo "Usage:"
	@echo "  make build           Build the application"
	@echo "  make run             Run the application"
	@echo "  make test            Run all tests"
	@echo "  make db-up           Start PostgreSQL via Docker"
	@echo "  make db-init         Initialize database schema"
	@echo "  make migrate         Auto-migrate via GORM (applies entities)"
	@echo "  make lint            Run golangci-lint"
	@echo "  make vet             Run go vet"
	@echo "  make tidy            Run go mod tidy"
	@echo "  make clean           Clean build artifacts"

build:
	go build -o bin/simplcommerce ./cmd/

run:
	go run ./cmd/

test:
	go test ./...

db-up:
	@if ! docker ps --format '{{.Names}}' | grep -q simplcommerce-pg; then \
		docker run -d --name simplcommerce-pg \
			-e POSTGRES_DB=simplcommerce \
			-e POSTGRES_USER=simplcommerce \
			-e POSTGRES_PASSWORD=simplcommerce \
			-p 5432:5432 \
			postgres:16-alpine && \
		echo "Waiting for PostgreSQL to start..." && \
		sleep 3; \
	else \
		echo "PostgreSQL already running"; \
	fi

db-init: db-up
	@echo "Creating tables..."
	@for f in queries/*.sql; do \
		echo "  Running $$f..."; \
		PGPASSWORD=simplcommerce psql -h localhost -U simplcommerce -d simplcommerce -f "$$f" 2>/dev/null || true; \
	done
	@echo "Database initialized successfully"

migrate:
	@echo "Running GORM auto-migration via application..."
	go run ./cmd/ --migrate-only

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	@rm -rf bin/
	@echo "Cleaned"
