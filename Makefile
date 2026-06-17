.PHONY: help build-all build-identity build-catalog build-cart build-orders build-payment build-shipping build-inventory build-pricing build-reviews build-cms build-tax build-activitylog build-notifications build-search build-gateway run-all run-identity run-catalog run-cart run-orders run-payment run-shipping run-inventory run-pricing run-reviews run-gateway run-gateway-dev test-all test-identity test-catalog clean docker-up docker-down db-init lint tidy vet

# Build all services
build-all:
	@echo "Building all services..."
	cd services/identity && go build ./cmd/identityd/
	cd services/catalog && go build ./cmd/catalogd/
	cd services/cart && go build ./cmd/cartd/
	cd services/orders && go build ./cmd/ordersd/
	cd services/payment && go build ./cmd/paymentd/
	cd services/shipping && go build ./cmd/shippingd/
	cd services/inventory && go build ./cmd/inventoryd/
	cd services/pricing && go build ./cmd/pricingd/
	cd services/reviews && go build ./cmd/reviewsd/
	cd services/cms && go build ./cmd/cmsd/
	cd services/tax && go build ./cmd/taxd/
	cd services/activitylog && go build ./cmd/activitylogd/
	cd services/notifications && go build ./cmd/notificationsd/
	cd services/search && go build ./cmd/searchd/
	cd services/api-gateway && go build ./cmd/gatewayd/
	@echo "All services built successfully"

# Individual service build targets
build-identity:
	cd services/identity && go build ./cmd/identityd/

build-catalog:
	cd services/catalog && go build ./cmd/catalogd/

build-cart:
	cd services/cart && go build ./cmd/cartd/

build-orders:
	cd services/orders && go build ./cmd/ordersd/

build-payment:
	cd services/payment && go build ./cmd/paymentd/

build-shipping:
	cd services/shipping && go build ./cmd/shippingd/

build-inventory:
	cd services/inventory && go build ./cmd/inventoryd/

build-pricing:
	cd services/pricing && go build ./cmd/pricingd/

build-reviews:
	cd services/reviews && go build ./cmd/reviewsd/

build-cms:
	cd services/cms && go build ./cmd/cmsd/

build-tax:
	cd services/tax && go build ./cmd/taxd/

build-activitylog:
	cd services/activitylog && go build ./cmd/activitylogd/

build-notifications:
	cd services/notifications && go build ./cmd/notificationsd/

build-search:
	cd services/search && go build ./cmd/searchd/

build-gateway:
	cd services/api-gateway && go build ./cmd/gatewayd/

# Run individual services
run-identity:
	cd services/identity && go run ./cmd/identityd/

run-catalog:
	cd services/catalog && go run ./cmd/catalogd/

run-cart:
	cd services/cart && go run ./cmd/cartd/

run-orders:
	cd services/orders && go run ./cmd/ordersd/

run-payment:
	cd services/payment && go run ./cmd/paymentd/

run-shipping:
	cd services/shipping && go run ./cmd/shippingd/

run-inventory:
	cd services/inventory && go run ./cmd/inventoryd/

run-pricing:
	cd services/pricing && go run ./cmd/pricingd/

run-reviews:
	cd services/reviews && go run ./cmd/reviewsd/

run-gateway:
	cd services/api-gateway && go run ./cmd/gatewayd/

# Run all services (in background)
run-all:
	@echo "Starting all services..."
	cd services/identity && go run ./cmd/identityd/ &
	cd services/catalog && go run ./cmd/catalogd/ &
	cd services/cart && go run ./cmd/cartd/ &
	cd services/orders && go run ./cmd/ordersd/ &
	cd services/payment && go run ./cmd/paymentd/ &
	cd services/shipping && go run ./cmd/shippingd/ &
	cd services/inventory && go run ./cmd/inventoryd/ &
	cd services/pricing && go run ./cmd/pricingd/ &
	cd services/reviews && go run ./cmd/reviewsd/ &
	cd services/api-gateway && go run ./cmd/gatewayd/ &
	@echo "All services started"

# Run tests for all services
test-all:
	cd pkg && go test ./...
	cd services/identity && go test ./...
	cd services/catalog && go test ./...
	cd services/cart && go test ./...
	cd services/orders && go test ./...
	cd services/payment && go test ./...
	cd services/shipping && go test ./...
	cd services/inventory && go test ./...
	cd services/pricing && go test ./...
	cd services/reviews && go test ./...
	cd services/cms && go test ./...
	cd services/tax && go test ./...
	cd services/activitylog && go test ./...
	cd services/notifications && go test ./...
	cd services/search && go test ./...
	cd services/api-gateway && go test ./...

# Test individual services
test-identity:
	cd services/identity && go test ./...

test-catalog:
	cd services/catalog && go test ./...

# Docker
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

# Database
db-init:
	psql -h localhost -U simplcommerce -d simplcommerce -f scripts/init-db.sql

# Lint
lint:
	golangci-lint run ./...

# Clean
clean:
	@echo "Cleaning build artifacts..."
	find . -type f -name "identityd" -delete
	find . -type f -name "catalogd" -delete
	find . -type f -name "cartd" -delete
	find . -type f -name "ordersd" -delete
	find . -type f -name "paymentd" -delete
	find . -type f -name "shippingd" -delete
	find . -type f -name "inventoryd" -delete
	find . -type f -name "pricingd" -delete
	find . -type f -name "reviewsd" -delete
	find . -type f -name "cmsd" -delete
	find . -type f -name "taxd" -delete
	find . -type f -name "activitylogd" -delete
	find . -type f -name "notificationsd" -delete
	find . -type f -name "searchd" -delete
	find . -type f -name "gatewayd" -delete

# Tidy all go modules
tidy:
	cd pkg && go mod tidy
	cd services/identity && go mod tidy
	cd services/catalog && go mod tidy
	cd services/cart && go mod tidy
	cd services/orders && go mod tidy
	cd services/payment && go mod tidy
	cd services/shipping && go mod tidy
	cd services/inventory && go mod tidy
	cd services/pricing && go mod tidy
	cd services/reviews && go mod tidy
	cd services/cms && go mod tidy
	cd services/tax && go mod tidy
	cd services/activitylog && go mod tidy
	cd services/notifications && go mod tidy
	cd services/search && go mod tidy
	cd services/api-gateway && go mod tidy

# Vet all modules
vet:
	cd pkg && go vet ./...
	cd services/identity && go vet ./...
	cd services/catalog && go vet ./...
	cd services/cart && go vet ./...
	cd services/orders && go vet ./...
	cd services/payment && go vet ./...
	cd services/shipping && go vet ./...
	cd services/inventory && go vet ./...
	cd services/pricing && go vet ./...
	cd services/reviews && go vet ./...
	cd services/cms && go vet ./...
	cd services/tax && go vet ./...
	cd services/activitylog && go vet ./...
	cd services/notifications && go vet ./...
	cd services/search && go vet ./...
	cd services/api-gateway && go vet ./...
