docker-deps:
	@echo "Starting dependencies..."
	@docker-compose up -d
go-tidy:
	@echo "Tidying up..."
	@cd go && go mod tidy
go-generate:
	@echo "Generating mocks..."
	@cd go && go generate ./...
go-test:
	@echo "Running tests..."
	@cd go && go test -v -cover -coverprofile=coverage.out ./...
	@cd go && go tool cover -html=coverage.out -o coverage.html
go-run-api:
	@set -a; \
	source .env; \
	set +a; \
	cd go && go run cmd/api/main.go