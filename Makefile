# Makefile

# Variables
BINARY_NAME=myapp
PROJECT_PATH=github.com/yourusername/myapp
DOCKER_IMAGE=myapp
API_DIR=api/graph

.PHONY: all setup generate build run test docker clean help

all: help

## help: Display list of commands
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## setup: Install dependencies
setup:
	@echo "Installing dependencies..."
	go get -u github.com/99designs/gqlgen@latest
	go get -u github.com/gin-gonic/gin
	go get -u go.mongodb.org/mongo-driver/mongo
	go mod tidy

## generate: Generate GraphQL code
generate:
	@echo "Generating GraphQL code..."
	gqlgen generate --config ${API_DIR}/gqlgen.yml --verbose

## run: Run the application locally
run:
	@echo "Starting development server..."
	go run cmd/main.go

## build: Build the binary
build:
	@echo "Building binary..."
	go build -o ${BINARY_NAME} cmd/main.go

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t ${DOCKER_IMAGE}:latest .

## docker-run: Run Docker container
docker-run:
	@echo "Starting Docker container..."
	docker-compose up --build

## docker-clean: Stop and remove Docker containers
docker-clean:
	@echo "Cleaning Docker environment..."
	docker-compose down -v
	docker rmi ${DOCKER_IMAGE}:latest || true

## clean: Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf ${API_DIR}/generated
	rm -rf ${API_DIR}/model
	rm -f ${API_DIR}/*.resolvers.go
	rm -f ${API_DIR}/generated.go

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Check code correctness
vet:
	@echo "Checking code..."
	go vet ./...

## migrate: Run database migrations
migrate:
	@echo "Running migrations..."
	# Add your migration commands here

## all: Full build pipeline
all: clean setup generate fmt vet build
