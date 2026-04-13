run-watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "'air' (hot reload) is not installed on your machine. Would you like to install it? [y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install it. Exiting..."; \
                exit 1; \
            fi; \
        fi

build:
	@echo "Building binary..."
	@go build -v -o bin/api cmd/api/main.go
	@echo "Build Finished!"


test:
	@echo "Running all tests..."
	@go test ./...

fmt:
	@golangci-lint fmt 

test-coverage:
	@echo "Running all tests and generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@echo "All tests finished!"

lint:
	@golangci-lint run

lint-fix:
	@golangci-lint run --fix

mig-new:
	@migrate create -ext sql -dir ./migrations $(name)

migrate-up:
	@migrate -path ./migrations -database "postgres://dev:dev@localhost:5432/dev?sslmode=disable&search_path=public" up

jet:
	@go run ./scripts/jet/jet.go

mig-up-jet: migrate-up jet

db-dev:
	@echo "Destroying any existing Postgres container..."
	@docker rm -f dev-postgres 2>/dev/null || true

	@echo "Starting fresh Postgres instance..."
	@docker run -d \
		--name dev-postgres \
		-e POSTGRES_USER=dev \
		-e POSTGRES_PASSWORD=dev \
		-e POSTGRES_DB=dev \
		-p 5432:5432 \
		postgres:17-alpine

	@echo "Waiting for Postgres to be ready..."

	@until docker exec dev-postgres pg_isready -U dev > /dev/null 2>&1; do \
		sleep 1; \
	done

	@echo "Postgres is ready to accept connections."
