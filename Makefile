.PHONY: run test lint tools migrate-up migrate-down migrate-create swagger mock

run:
	AUTH_PORT=8080 \
	DATABASE_URL=postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable \
	REDIS_URL=redis://localhost:6379/0 \
	go run ./cmd/server

test:
	go test ./... -v -count=1

lint:
	golangci-lint run --timeout=5m

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/vektra/mockery/v2@latest

mock:
	mockery

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-create:
	goose -dir migrations create $(NAME) sql

swagger:
	swag init --parseDependency -g cmd/server/main.go -o docs/
