.PHONY: help build run test clean run-db lint vet tidy

help:
	@echo "SimplCommerce-Go"
	@echo ""
	@echo "Usage:"
	@echo "  make build        Build the application"
	@echo "  make run          Run the application"
	@echo "  make run-db       Run PostgreSQL via Docker"
	@echo "  make test         Run all tests"
	@echo "  make lint         Run golangci-lint"
	@echo "  make vet          Run go vet"
	@echo "  make tidy         Run go mod tidy"
	@echo "  make clean        Clean build artifacts"

build:
	go build -o bin/simplcommerce ./cmd/

run:
	go run ./cmd/

test:
	go test ./...

run-db:
	docker run -d --name simplcommerce-pg \
		-e POSTGRES_DB=simplcommerce \
		-e POSTGRES_USER=simplcommerce \
		-e POSTGRES_PASSWORD=simplcommerce \
		-p 5432:5432 \
		postgres:16-alpine

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	@rm -rf bin/
	@echo "Cleaned"
