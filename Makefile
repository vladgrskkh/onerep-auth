.PHONY: run test lint migrate-up migrate-down migrate-create

run:
	AUTH_PORT=8080 \
	DATABASE_URL=postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable \
	REDIS_URL=redis://localhost:6379/0 \
	go run ./cmd/server

test:
	go test ./... -v -count=1

lint:
	golangci-lint run --timeout=5m

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-create:
	goose -dir migrations create $(NAME) sql
