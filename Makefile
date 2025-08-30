.DEFAULT_GOAL := help

# Variables
APP_NAME=app
MAIN_FILE=cmd/server/main.go
WORKER_FILE=cmd/worker/consumer/main.go
DOCKER_REGISTRY=ghcr.io
DOCKER_REGISTRY_APP=fiap-soat-g20/tc4-order-service
DOCKER_REGISTRY_MOCK_SERVER_APP=fiap-soat-g20/mock-server
VERSION=$(shell git describe --tags --always --dirty)
NAMESPACE=tech-challenge-ns
TEST_PATH=./internal/...
TEST_COVERAGE_FILE_NAME=coverage.out
MIGRATION_PATH = internal/infrastructure/database/migrations
DB_URL = postgres://postgres:postgres@localhost:5432/fastfood_10soat_g19_tc4_order?sslmode=disable

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOTIDY=$(GOCMD) mod tidy

# Looks at comments using ## on targets and uses them to produce a help output.
.PHONY: help
help: ALIGN=22
help: ## Print this message
	@echo "Usage: make <command>"
	@awk -F '::? .*## ' -- "/^[^':]+::? .*## /"' { printf "  make '$$(tput bold)'%-$(ALIGN)s'$$(tput sgr0)' - %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Format the code
	@echo  "🟢 Formatting the code..."
	$(GOCMD) fmt ./...

.PHONY: build
build: fmt ## Build the application
	@echo  "🟢 Building the application..."
	$(GOBUILD) -o bin/$(APP_NAME) $(MAIN_FILE)

.PHONY: run-db
run-db: ## Run the database
	@echo  "🟢 Running the database..."
	docker-compose up -d db dbadmin

.PHONY: run-api
run-api: build run-db ## Run the API application
	@echo  "🟢 Running the application..."
	$(GORUN) $(MAIN_FILE) || true


.PHONY: run-worker
run-worker: build run-db ## Run the worker application 
	@echo  "🟢 Running the application..."
	$(GORUN) $(WORKER_FILE) || true

.PHONY: stop
stop: ## Stop the application
	@echo  "🔴 Stopping the application..."
	docker-compose down	

.PHONY: stop-db
stop-db: ## Stop the database
	@echo  "🔴 Stopping the database..."
	docker-compose down db dbadmin

.PHONY: run-api-air
run-api-air: build ## Run the application with Air
	@echo  "🟢 Running the application with Air..."
	@go tool air -c air.toml

.PHONY: tests
tests: lint ## Run tests
	@echo  "🟢 Running tests..."
	@$(GOFMT) ./...
	@$(GOVET) ./...
	@$(GOTIDY)
	$(GOTEST) $(TEST_PATH) -race -cover

.PHONY: coverage
coverage: ## Run tests with coverage
	@echo  "🟢 Running tests with coverage..."
# remove files that are not meant to be tested
	$(GOTEST) $(TEST_PATH) -race -cover -coverprofile=$(TEST_COVERAGE_FILE_NAME).tmp
	@cat $(TEST_COVERAGE_FILE_NAME).tmp | grep -v "_mock.go" | grep -v "_request.go" | grep -v "_response.go" \
	| grep -v "_gateway.go" | grep -v "_datasource.go" | grep -v "_presenter.go" | grep -v "middleware" \
	| grep -v "config" | grep -v "route" | grep -v "util" | grep -v "database" \
	| grep -v "server" | grep -v "logger" | grep -v "httpclient" | grep -v "_entity.go" | grep -v "errors.go" | grep -v "_dto.go" > $(TEST_COVERAGE_FILE_NAME)
	@rm $(TEST_COVERAGE_FILE_NAME).tmp
	$(GOCMD) tool cover -html=$(TEST_COVERAGE_FILE_NAME)

.PHONY: clean
clean: ## Clean up binaries and coverage files
	@echo "🔴 Cleaning up..."
	$(GOCLEAN)
	rm -f $(APP_NAME)
	rm -f $(TEST_COVERAGE_FILE_NAME)

.PHONY: mock
mock: ## Generate mocks
	@echo  "🟢 Generating mocks..."
# romove mocks files
	@rm -rf internal/core/port/mocks/*
# loop through all files in the port directory and generate mocks
	@for file in internal/core/port/*.go; do \
		go tool mockgen -source=$$file -destination=internal/core/port/mocks/`basename $$file _port.go`_mock.go; \
	done

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo  "🟢 Generating Swagger documentation..."
	@go tool swag fmt ./...
	@go tool swag init -g ${MAIN_FILE} --parseInternal true

.PHONY: lint
lint: ## Run linter
	@echo  "🟢 Running linter..."
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2 run

.PHONY: migrate-create
migrate-create: ## Create new migration, usage example: make migrate-create name=create_table_products
	@echo  "🟢 Creating new migration..."
# if name is not passed, required argument
ifndef name
	$(error name is not set, usage example: make migrate-create name=create_table_products)
endif
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2 create -ext sql -dir ${MIGRATION_PATH} -seq $(name)

.PHONY: migrate-up
migrate-up: ## Run migrations
	@echo  "🟢 Running migrations..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2 -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose up

.PHONY: migrate-down
migrate-down: ## Roll back migrations
	@echo  "🔴 Rolling back migrations..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2 -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose down

.PHONY: install
install: ## Install dependencies
	@echo  "🟢 Installing dependencies..."
	go mod download

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo  "🟢 Building Docker image..."
	docker build --platform linux/amd64 -t $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):$(VERSION) .
	docker tag $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):$(VERSION) $(DOCKER_REGISTRY)/$(APP_NAME):latest

.PHONY: docker-push
docker-push: ## Push Docker image
	@echo  "🟢 Pushing Docker image..."
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):latest

.PHONY: compose-build
compose-build: ## Build the application with Docker Compose
	@echo "🟢 Building the application..."
	docker compose build

.PHONY: compose-up
compose-up: ## Start development environment with Docker Compose
	@echo  "🟢 Starting development environment..."
	docker compose pull
	docker-compose up -d --wait --build

.PHONY: compose-down
compose-down: ## Stop development environment with Docker Compose
	@echo  "🔴 Stopping development environment..."
	docker-compose down

.PHONY: compose-clean
compose-clean: ## Clean the application with Docker Compose, removing volumes and images
	@echo "🔴 Cleaning the application..."
	docker compose down --volumes --rmi all

.PHONY: scan
scan: ## Run security scan
	@echo  "🟢 Running security scan..."
	@go tool govulncheck -show verbose ./...
	@go tool trivy image --severity HIGH,CRITICAL $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):latest

.PHONY: new-branch
new-branch: ## Create new branch
	@echo "🟢 Creating new branch..."
	./scripts/new-branch.sh -c

.PHONY: pull-request
pull-request: ## Create pull request
	@echo "🟢 Creating pull request..."
	./scripts/pull-request.sh

.PHONY: bdd-tests
bdd-tests: ## Run BDD tests
	@echo "🟢 Running BDD tests..."
	go test -test.v -test.run ^TestFeatures$$ ./tests

