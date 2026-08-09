# onerep-auth

Authentication and user profile service for OneRep gym training app.

> **Follow [`CONVENTIONS.md`](https://github.com/vladgrskkh/gym/blob/main/CONVENTIONS.md) — the shared engineering conventions for all OneRep services. This file lists the service-specific parts.**

## Stack
- Go 1.26+, PostgreSQL 16, Redis 7

## Commands
```bash
make run          # Start server on :8080
make test         # Run all tests
make lint         # golangci-lint run (golden config)
make tools        # Install pinned tools
make generate     # mockery + swagger
make mock         # mockery only
make swagger      # swag init only
make check-generate  # regenerate + fail on diff
make migrate-up   # Apply migrations
make migrate-down # Rollback migrations
```

## Architecture
Domain-based layered structure (see CONVENTIONS.md):

```
cmd/server/              # Entry point (thin)
internal/
  domain/
    auth/                # User, OAuthAccount, Gender, errors — pure, no JSON
    user/                # UserProfile read model, UpdateProfileInput, ToProfile
  service/
    auth/                # AuthService, OAuthService + consumer interfaces
    user/                # UserService (operates on domain types)
  handler/
    auth/                # handlers + dto/ + error_mapper.go
    user/                # handlers + dto/ + error_mapper.go + mapper.go
  infrastructure/
    auth/
      postgres/          # UserRepo, OAuthAccountRepo
      redis/             # TokenStore
      jwt/               # TokenManager, Claims
      crypto/            # PasswordHasher
      oauth/             # Google adapter
  application/           # App wiring (options pattern), route registration
  config/                # caarlos0/env
  handler/               # shared HTTP utils: response, context, middleware
```

## Service-specific notes
- JWT: RS256 access tokens (15m) + opaque refresh tokens (7d, Redis, rotation).
- JWKS endpoint at `/.well-known/jwks.json` for cross-service validation.
- Email/password + Google OAuth (Apple pending).
