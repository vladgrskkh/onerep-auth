# onerep-auth

Authentication and user profile service for OneRep gym training app.

## Stack
- Go 1.23+, PostgreSQL 16, Redis 7

## Commands
```bash
make run          # Start server on :8080
make test         # Run all tests
make lint         # golangci-lint run
make migrate-up   # Apply migrations
make migrate-down # Rollback migrations
```

## Architecture
Domain-based structure with CSP (Consumer-defined interfaces):

```
cmd/server/          # Entry point
internal/
  auth/              # Auth domain: User, AuthService, OAuthService, handlers, consumer interfaces
  user/              # User domain: UserService, UpdateProfileInput, handlers
  domain/            # Shared error sentinels (no imports)
  application/       # App wire-up, route registration
  config/            # Config with caarlos0/env
  handler/           # Shared HTTP utilities: response helpers, context, error mapping, DTOs, middleware
  infrastructure/    # Postgres, Redis, JWT, Crypto, OAuth implementations
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
Interfaces are declared where they are consumed (auth/, user/), not in infrastructure/.
No central interfaces.go file. Each service declares only the methods it needs.

## Mock generation
```bash
mockery
```
