# onerep-auth

Authentication and user profile service for OneRep gym training app.

## Stack
- Go 1.23+, PostgreSQL 16, Redis 7

## Commands
```bash
make run          # Start server on :8080
make test         # Run all tests
make lint         # go vet ./...
make migrate-up   # Apply migrations
make migrate-down # Rollback migrations
```

## Architecture
DDD with CSP (Consumer-defined interfaces):

```
cmd/server/          # Entry point
internal/
  domain/            # Entities, value objects — no interfaces, no imports
  application/       # Use cases + consumer-defined interfaces (ISP)
  handler/           # HTTP handlers, middleware
  infrastructure/    # Postgres, Redis, JWT, OAuth implementations
```

## CI
- `golangci-lint` v2 with golden config (maratori)
- Docker build pushes to `ghcr.io/vladgrskkh/onerep-auth`

## Local dev
```bash
# Start dependencies
docker compose -f ../docker-compose.yml up -d

# Run
make run
```

## ISP rule
Interfaces are declared where they are consumed (application/), not in infrastructure/.
No central interfaces.go file. Each service declares only the methods it needs.
