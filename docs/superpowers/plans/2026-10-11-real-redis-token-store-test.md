# Real Redis Token Store Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the hand-written token-store fake test with a real Redis integration regression test that proves one-time consumption through Redis `GETDEL`.

**Architecture:** Keep the production `TokenStore` and its `redisClient` abstraction unchanged. The adapter test will parse `REDIS_URL` (defaulting to `redis://localhost:6379/0`), create the real `go-redis/v9` client, ping Redis, skip clearly when unavailable, and clean up a uniquely namespaced token.

**Tech Stack:** Go 1.26.5, `github.com/redis/go-redis/v9`, Redis 7, `testify/suite`.

**Spec:** User request in the current conversation.

## Global Constraints

- Do not change production behavior.
- Use the real `github.com/redis/go-redis/v9` client in the adapter regression test.
- Read `REDIS_URL`, defaulting to `redis://localhost:6379/0`.
- Ping Redis and skip with a clear reason when no local Redis is available.
- Save a uniquely namespaced token and clean it up.
- Assert the first `Consume` returns the user ID and the second returns the repository’s not-found error.
- Keep the existing service concurrency test.
- Run the requested focused and repository-wide checks, then commit without pushing or opening a PR.

## Review Focus

- Missing local Redis: the test skips with the configured URL and ping error instead of failing normal unit runs — adapter test.
- A malformed `REDIS_URL`: setup fails explicitly rather than silently using another server — adapter test setup.
- Token reuse: the first consume returns the saved user ID and the second returns `token not found` — adapter test.
- Cross-test key collisions: the saved token includes a unique test namespace and cleanup deletes its key — adapter test.
- Existing service atomicity coverage: both concurrent refresh calls still produce exactly one not-found result — existing service test remains unchanged.

---

### Task 1: Replace the fake adapter test with real Redis coverage

**Files:**
- Modify: `internal/infrastructure/auth/redis/token_store_test.go`
- Inspect only: `internal/infrastructure/auth/redis/token_store.go`, `go.mod`, `Makefile`, `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `authredis.NewTokenStore`, `Save`, `Consume`, and the existing repository error text.
- Produces: a test-only real Redis setup using `redislib.ParseURL`, `redislib.NewClient`, `Ping`, and cleanup via `Del`.

- [ ] **Step 1: Write the failing real-client test first**

  Remove the fake client and make the test construct a real `*redislib.Client` from `REDIS_URL` or `redis://localhost:6379/0`. Ping with a bounded context and skip clearly on connection failure. Save a token with a unique namespace, register cleanup, then assert the two `Consume` calls return the user ID and `token not found` respectively.

- [ ] **Step 2: Run the focused test to verify the red/skip result**

  Run `go test ./internal/infrastructure/auth/redis -run TestTokenStoreSuite -count=1 -v` and record either the expected failure from the incomplete test setup or the explicit Redis-unavailable skip.

- [ ] **Step 3: Finish the test setup minimally**

  Keep only the real client, bounded ping, unique token, cleanup, and behavior assertions required by the regression. Do not modify `token_store.go` or the service concurrency test.

- [ ] **Step 4: Run the focused test to verify the real regression**

  Run `go test ./internal/infrastructure/auth/redis -run TestTokenStoreSuite -count=1 -v`; with Redis 7 available it must pass and exercise actual `GETDEL`, otherwise it must skip with the configured endpoint and reason.

### Task 2: Verify the repository and commit

**Files:**
- Inspect: all repository files affected by the requested checks.

- [ ] **Step 1: Run focused tests and all requested checks**

  Run focused tests, `go test ./... -count=1`, `go test -race ./...`, `go vet ./...`, `make lint`, and `make check-generate` where the environment permits. Capture exit status and any environment skips/failures exactly.

- [ ] **Step 2: Review the diff and preserve required coverage**

  Confirm production files are unchanged, the fake adapter is gone, CI still uses Redis 7, and `TestRefresh_ConcurrentRequestsConsumeTokenOnlyOnce` remains present and unchanged.

- [ ] **Step 3: Commit**

  Commit the verified test-only change with an informative conventional message such as `test(auth): cover token store with real redis`.
