# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A Go REST API for movie-theatre booking: chi router, PostgreSQL via pgx/v5, sqlc-generated queries, goose migrations, JWT auth. Module path: `github.com/Okwujohngibson25/movie-go`.

## Commands

```bash
air                                  # run with hot reload (builds ./cmd/api into ./tmp/main, per .air.toml)
go run ./cmd/api                     # run without hot reload
go build ./...
go test ./...
go test ./internal/helper -run TestDecodeJSON_RejectsUnknownFieldAndOversize   # single test
sqlc generate                        # regenerate internal/database/db from db/queries + migrations
goose -dir migrations postgres "$DATABASE_URL" up    # apply migrations (goose is a CLI, not a go.mod dep)
```

Config is loaded from `.env` by `godotenv` at startup; the app fails to start if `.env` is missing. Keys: `PORT` (passed straight to `http.ListenAndServe`, so it needs the colon, e.g. `:8080`), `APP_ENV`, `JWT_SECRET`, `DATABASE_URL`, `Admin_email`, `Admin_Password` (note the mixed casing). An admin user is seeded from the last two on every startup if it doesn't already exist.

`api/**/*.http` files are hand-written request samples for exercising endpoints.

## Architecture

Each feature (`internal/movie`, `screen`, `showing`, `user`, and the still-empty `booking`) follows the same layering, wired together manually in `cmd/api/main.go`:

- `handler.go`: decodes the request with `helper.DecodeJSON` (1 MiB cap, rejects unknown fields and trailing JSON), validates the `dto.go` request struct with the shared validator, maps it to an `input.go` struct, calls the service, and writes `response.APIResponse` via `response.WriteJSON`. `mapper.go` turns domain models into response DTOs.
- `service.go`: business logic. It **defines its own repository interface** (consumer-side), working only in `internal/domains` types. Services that need another feature's data declare a narrow interface for it (e.g. `showing.Movie` is satisfied by the movie repo, which is used to compute `EndTime` from the movie's duration).
- `internal/postgres/*.go`: repository implementations wrapping the sqlc `db.Queries`. They convert between `uuid.UUID`/`decimal.Decimal` and `pgtype.UUID`/`pgtype.Numeric` (see `helper.NumericToDecimal`), and translate `pgx.ErrNoRows` or zero affected rows into `apperrors` sentinels.

Data layer:

- Schema lives in `migrations/` (goose `-- +goose Up/Down` files); sqlc reads it as the schema source. Queries live in `db/queries/*.sql`. **`internal/database/db/` is generated; don't edit it by hand**, change the SQL and run `sqlc generate`.
- Movies are soft-deleted (`deleted_at`), and queries filter on `deleted_at IS NULL`. `CreateMovie` uses `ON CONFLICT (title) DO NOTHING RETURNING *`, so a duplicate shows up as `pgx.ErrNoRows` and becomes `ErrMovieAlreadyExists`.

Errors: return sentinel errors from `internal/apperrors` (wrapping with `%w` is fine). `apperrors.HandleError` is the single place that maps errors (including `validator.ValidationErrors`) to HTTP status codes and response `code` strings; anything unmapped is logged and returned as a 500. A new error type needs both a sentinel and a `case` there. (`ErrMovieAlreadyExists` currently has no case, so it falls through to 500.)

Validation: `internal/validation.NewValidator` builds one shared `*validator.Validate` and registers custom tags (e.g. `rating`, a decimal between 1 and 5).

Auth:

- `POST /api/v1/auth/register` and `/login` are public. Login issues a 15-minute access token and a 7-day refresh token (HS256, `JWT_SECRET`).
- Everything under `/api/v1/admin` goes through `middleware.RequireAuth` and then `RequireAdmin`. `RequireAuth` reads the **raw token from the `Authorization` header, with no `Bearer ` prefix**.
- The claims struct is defined twice: `user.CustomClaims` (signing, no JSON tags) and `middleware.CustomClaims` (parsing, tagged). They only line up because Go's JSON decoding is case-insensitive, so keep the field names in sync.

## Known quirks

- Movie request DTOs use the JSON key `duration_mintes` (typo), while responses use `duration_minutes`. Because `DecodeJSON` rejects unknown fields, clients have to send the misspelled key.

# Working rules for this repo

## Task workflow (strict)

- The task list lives in `docs/TASKS.md`. Progress notes live in `docs/PROGRESS.md`.
- Work on exactly ONE task at a time: the first task whose status is `todo`.
- Never start the next task until I explicitly say "approved".
- Follow test-first for every task:
  1. Restate the task's requirements as a numbered list in PROGRESS.md.
  2. Write tests that encode those requirements. Run them and show they FAIL.
  3. STOP. Set status to `tests-review` and wait for my approval of the tests.
  4. After approval, write the minimum code to make the tests pass.
     Do NOT modify approved tests. If a test seems wrong, stop and ask.
  5. Run: `go build ./... && go vet ./... && go test ./...`. All must pass.
  6. Update PROGRESS.md (what changed, files touched, open questions).
  7. STOP. Set status to `code-review` and wait.
- On "approved" after code review: set status to `done` and commit with message
  `task <id>: <title>`. Then stop. Do not begin the next task in the same turn.

## Project conventions

- Go, Postgres. Repository tests run against real Postgres (see section 5 task).
- Run `gofmt` on changed files.

0. Before starting a `todo` task, create and switch to a branch named
   `<id>-<short-title>` (e.g. `t0-refresh-token`) from `main`.

Use the Read, Glob and Grep tools to inspect files. Do not use Bash commands (cat, loops, find) to read files.
