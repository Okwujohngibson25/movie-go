# Go Skills Assessment — movie-go

_Reviewed: 2026-10-23 · Focus: Go conventions, Go patterns, and backend-system completeness (application logic not graded)_

## Grades

| | Score |
|---|---|
| **Go conventions & patterns** | **72 / 100** |
| **Backend-system completeness** | **35 / 100** (it's a solid core; most of the surrounding system isn't built yet) |

### Conventions & patterns breakdown

| Area | Weight | Score |
|---|---|---|
| Project layout & package design | 15 | 12 |
| Naming conventions | 15 | 8 |
| Interfaces & dependency injection | 15 | 12 |
| Error-handling patterns | 15 | 12 |
| Concurrency/context usage | 10 | 8 |
| Type design (structs, constants, generics) | 10 | 6 |
| Testing patterns | 10 | 8 |
| Code economy (no redundancy, idiomatic brevity) | 10 | 6 |
| **Total** | **100** | **72** |

---

## Part 1 — Go conventions & patterns

### What you're doing well

**Project layout**
- `cmd/api` for the binary and `internal/` for everything else is the standard layout. `internal/` stops other modules importing your packages.
- **Package-by-feature** (`movie`, `screen`, `showing`, `user`) rather than package-by-layer (`handlers/`, `services/`). Go favours this approach.
- Generated code is isolated in `internal/database/db`, and the SQL sits in `db/queries`.
- The URL path is versioned (`/api/v1`).

**Interfaces & DI**
- **Accept interfaces, return structs.** Constructors return concrete `*MovieService`, but services accept a `MovieRepository` interface.
- **Interfaces are defined by the consumer**, not by the implementer. This is the single most important Go interface convention, and you've got it right.
- **Small interfaces.** `showing.Movie` has one method, in the spirit of `io.Reader`.
- `NewX(...)` constructors and manual wiring in `main`, with no DI container and no package-level globals.

**Errors**
- Sentinel errors are named `ErrXxx` with lowercase, unpunctuated messages, which matches the Go Code Review Comments.
- Wrapping uses `%w`, and inspection uses `errors.Is` / `errors.As`. You never compare error strings.
- Each layer translates driver errors (`pgx.ErrNoRows`, `*pgconn.PgError`) into domain errors at its boundary.
- Handlers use the early-return guard-clause style throughout.

**Context & resources**
- `ctx context.Context` is always the first parameter and is propagated from `r.Context()`.
- `defer tx.Rollback(ctx)` followed by `tx.Commit(ctx)` is the canonical transaction pattern.
- `defer pool.Close()` in `main`.
- An unexported `type contextKey string` is used for context values, which avoids key collisions.

**HTTP patterns**
- Middleware has the `func(http.Handler) http.Handler` shape, with a closure for configuration (`RequireAuth(secret)`).
- chi `Route` and `Group` are used for scoping middleware.
- `http.MaxBytesReader` and `json.Decoder.DisallowUnknownFields()` are used for input.

**Testing patterns**
- Table-driven tests with `t.Run` subtests.
- `t.Helper()` on helpers and `t.Cleanup()` to restore global state.
- Hand-written fakes instead of a mocking framework (the Go-preferred style).
- Tests sit next to the code in the same package, so they can test unexported functions like `movieFromRow`.

---

### Where your conventions could be better

#### 1. Naming

| Issue | Example | Convention |
|---|---|---|
| **Stutter** | `movie.MovieService`, `movie.NewMovieHandler`, `postgres.MovieSqlcRepository` | The package name is part of the identifier. Use `movie.Service`, `movie.NewHandler`, `postgres.MovieRepo`. |
| **Inconsistent receivers** | Handlers use `s`, services use `r` (which collides with `r *http.Request`) | Use a short, consistent abbreviation of the type: `h *Handler`, `s *Service`, `r *MovieRepo`. |
| **Underscores** | `User_ID` | Use MixedCaps only: `UserID`. |
| **Initialisms** | `GetScreensByID(ctx, Id uuid.UUID)`, `movieIdStr` | Write initialisms in a single case: `id`, `ID`, `movieIDStr`. |
| **Capitalised locals** | `PasswordString`, `ScreenData`, `CreateUserParams`, `ValidationErrors` | Locals are lowerCamel. `ValidationErrors` also shadows the type name. |
| **Plural type names** | `domains.Seats` (a single seat) | Type names are singular: `Seat`. |
| **Plural/grab-bag packages** | `domains`, `helper` | Package names are singular nouns that describe their contents (`domain`). Avoid `helper`/`util`/`common`; split by purpose, e.g. `httpx` (DecodeJSON) and `pgconv` (NumericToDecimal, UUID conversion). |
| **Odd named returns** | `(signedAccess string, rawRefresh string, errordata error)` | Use named returns only when they document something; never name the error anything but `err`. |

#### 2. Handlers depend on concrete types
You applied consumer-side interfaces between services and repositories but not between handlers and services. `MovieHandler` holds `*MovieService`, so handlers can't be unit-tested without the whole stack. Use the same pattern one layer up:

```go
// movie/handler.go
type service interface {
    Create(ctx context.Context, in Input) error
    Get(ctx context.Context, id uuid.UUID) (domain.Movie, error)
}

type Handler struct {
    svc      service
    validate *validator.Validate
}
```

#### 3. Compile-time interface checks
Add a guard in each repository file so a signature drift fails at compile time, not at the wiring in `main`:

```go
var _ movie.Repository = (*MovieRepo)(nil)
```

#### 4. Redundant error boilerplate
About 15 functions follow this shape:
```go
if err := r.repo.DeleteMovieByID(ctx, id); err != nil {
    return err
}
return nil
```
The idiomatic form is `return r.repo.DeleteMovieByID(ctx, id)`. Only add an `if err != nil` when you're adding context (`fmt.Errorf("delete movie %s: %w", id, err)`), which is the other half of the habit. Right now most errors bubble up without context about where they came from.

#### 5. Error type pattern vs. a giant switch
`HandleError` is a 100-line `switch`. In Go, the scalable pattern is an error type that carries its own HTTP semantics:

```go
type Error struct {
    Status int
    Code   string
    Msg    string
}

func (e *Error) Error() string { return e.Msg }

var ErrMovieNotFound = &Error{http.StatusNotFound, "MOVIE_NOT_FOUND", "movie not found"}

// HandleError
var appErr *apperrors.Error
if errors.As(err, &appErr) {
    response.WriteJSON(w, appErr.Status, ...)
    return
}
```
With this, adding an error no longer requires editing a central switch.

#### 6. Magic strings → typed constants
`"admin"`, `"user"`, `"access"` and `"refresh"` are scattered string literals. Go uses named types with constants for this:

```go
type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
)
```
The same applies to token types and the `bcrypt` cost `10` (use `bcrypt.DefaultCost` or a named constant).

#### 7. Type design
- **Don't duplicate types across packages.** `user.CustomClaims` and `middleware.CustomClaims` describe the same thing. Put one `auth.Claims` type in its own package and use it for both signing and parsing.
- **Generics that add nothing.** `APIResponse[T any]` is only ever `APIResponse[any]`. Either use the type parameter for real (`APIResponse[[]MovieDTO]`) or remove it. Don't add generics speculatively.
- **Domain types shouldn't carry `json` tags** (`domain.ShowingMovie`), and response DTOs shouldn't embed domain types. Keep the transport and domain layers separate, as you did for movies.
- **One request type is enough when two are identical.** `CreateMovieRequest` and `UpdateMovieRequest` are identical; one type or a shared embedded struct is enough until they diverge.
- **Pick integer types once.** `Capacity` goes `int32 → int → int32`, and `SeatNumber` goes `int32 → int64`. Choose one at the domain level and convert only at the database boundary.

#### 8. Repeated conversion code
`pgtype.UUID{Bytes: id, Valid: true}` appears around 15 times. A two-line `func pgUUID(id uuid.UUID) pgtype.UUID` removes the noise. Alternatively, configure sqlc `overrides` to generate `uuid.UUID` and `decimal.Decimal` directly, which removes most of the conversion layer.

#### 9. Ignored return values
`response.WriteJSON` returns an `error` that every caller discards. Either handle it inside `WriteJSON` and return nothing, or check it at the call site. A linter (`errcheck`) flags this.

#### 10. `main` structure
Use the `run() error` pattern. It keeps `main` to three lines, makes `defer`s run (`log.Fatal` skips them), and makes startup testable:

```go
func main() {
    if err := run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```
Also: `SeedAdmin` takes a `cfg` the service already holds, creates its own `context.Background()`, and returns nothing. Pass `ctx` in and return an `error`.

#### 11. Doc comments & tooling
- Exported identifiers have no doc comments (`// Service handles movie business logic.`). `go doc` and `revive` expect them.
- Commented-out code should be deleted, because git keeps history.
- Adopt **`golangci-lint`** with `errcheck`, `revive`, `staticcheck`, `gocritic` and `unparam`. It will flag most of the items above automatically.

#### 12. Testing patterns to add
- `httptest.NewRecorder()` handler tests (possible once item 2 is done).
- `go test -race ./...` as a habit.
- Integration tests against real Postgres using `testcontainers-go`, or a `TEST_DATABASE_URL` env var with `t.Skip` when it's unset.
- `t.Parallel()` on independent subtests.

---

## Part 2 — What's missing for a complete backend system

You said you know logging is missing, so it isn't repeated here. Everything else, grouped by concern:

### Server lifecycle
- [ ] `http.Server` with `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` (bare `ListenAndServe` has none)
- [ ] **Graceful shutdown**: `signal.NotifyContext` plus `srv.Shutdown(ctx)`, closing the DB pool after in-flight requests finish
- [ ] Per-request timeout middleware (`chi/middleware.Timeout`)
- [ ] DB pool tuning (`pgxpool.Config`: `MaxConns`, `MaxConnIdleTime`, `HealthCheckPeriod`)

### Cross-cutting middleware
- [ ] Panic recovery (`middleware.Recoverer`)
- [ ] Request ID / correlation ID (`middleware.RequestID`)
- [ ] CORS (`go-chi/cors`) if a browser frontend will call it
- [ ] Rate limiting, especially on `/auth/login` and `/auth/register` (`httprate`)
- [ ] Security headers (HSTS, `X-Content-Type-Options`, etc.)

### Health & observability (beyond logs)
- [ ] `GET /healthz` (liveness) and `GET /readyz` (pings the DB)
- [ ] Metrics: Prometheus `/metrics` with request count, latency and status by route
- [ ] Distributed tracing: OpenTelemetry (`otelhttp` and the pgx tracer)

### Auth & identity
- [ ] `POST /auth/refresh` (the refresh token is issued but can't be used)
- [ ] `POST /auth/logout` with server-side refresh-token storage, rotation and revocation
- [ ] `GET /me` / profile endpoints
- [ ] Password reset and email verification flows
- [ ] Non-admin, user-facing routes (browse movies and showings, make or view own bookings)
- [ ] `Authorization: Bearer <token>` header format (the standard)

### API design
- [ ] **Pagination** on every list endpoint (`?limit=&cursor=` or `?page=&size=`)
- [ ] Filtering and sorting (showings by date, movie or screen)
- [ ] Consistent resource responses (create should return the created resource or its `Location`, not just a message)
- [ ] **OpenAPI/Swagger spec** (hand-written, or generated with `swaggo` or `oapi-codegen`)
- [ ] Idempotency keys for non-idempotent operations (booking, payment)

### Domain features a booking system needs
- [ ] Booking service (still empty)
- [ ] Seat holds or reservations with expiry
- [ ] Payments integration and webhooks
- [ ] Notifications (email/SMS confirmation)
- [ ] **Background workers**: release expired holds, send emails (a goroutine worker with `context` cancellation, or a job queue such as `river`, which is Postgres-backed and Go-native)

### Data layer
- [ ] `created_at` / `updated_at` columns on all tables
- [ ] Audit trail for admin actions
- [ ] Indexes on foreign keys and query filters
- [ ] Caching for hot reads (movie and showing lists), in memory or Redis
- [ ] Migrations run as part of deploy (goose as a library in a `migrate` subcommand, or a CI step)
- [ ] Seed data for local development

### Configuration
- [ ] Validate required config at startup (fail fast on an empty `JWT_SECRET`)
- [ ] Don't require a `.env` file; real deployments inject environment variables
- [ ] Struct-tag config loading (`caarlos0/env` or `kelseyhightower/envconfig`) with defaults
- [ ] Consistent env-var casing (`ADMIN_EMAIL`, not `Admin_email`)

### Delivery & tooling
- [ ] `Dockerfile` (multi-stage build, distroless or scratch final image)
- [ ] `docker-compose.yml` for Postgres plus the API for local development
- [ ] `Makefile` or `Taskfile` (`make run`, `make test`, `make lint`, `make migrate`, `make sqlc`)
- [ ] CI pipeline (GitHub Actions: `build`, `vet`, `golangci-lint`, `test -race`, `sqlc diff`)
- [ ] `README.md` with setup, environment variables and API overview
- [ ] Integration and end-to-end tests in CI against a real database

---

## Suggested next steps (in order)

1. Install `golangci-lint` and fix what it reports. That covers most of Part 1 automatically.
2. Rename to remove stutter, fix receivers, and add handler-side interfaces and compile-time checks.
3. Introduce the `apperrors.Error` type and typed `Role` / token-type constants.
4. Add the `run()` pattern, `http.Server` timeouts, graceful shutdown, health checks and recovery middleware. Together they make about an afternoon of work.
5. Add a Dockerfile, docker-compose, a Makefile and CI.
6. Then build out the missing features: refresh/logout, pagination, OpenAPI, and the booking flow.

Steps 1 to 3 would move the conventions grade into the **mid-80s**.
