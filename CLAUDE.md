# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Bidgo is a Go HTTP API (auction/bidding service, early scaffold) using chi for routing, PostgreSQL via pgx, tern for migrations, and sqlc for generated query code. Module path: `github.com/lucaasscm/bidgo`.

## Commands

```bash
docker compose up -d          # start Postgres (reads .env; copy .env.example first)
air                           # dev server with hot reload (builds ./cmd/api into ./bin/api)
go run ./cmd/api              # run the API without hot reload
go build -o ./bin/api ./cmd/api

go run ./cmd/terndotenv       # apply database migrations (see note below)
tern new -m ./internal/store/pgstore/migrations <name>   # create a new migration

sqlc generate -f ./internal/store/pgstore/sqlc.yml       # regenerate query code

go test ./...                 # run tests
go test ./internal/api -run TestName   # run a single test
```

**Migrations must be run through `go run ./cmd/terndotenv`, not `tern migrate` directly.** The tern config (`internal/store/pgstore/migrations/tern.conf`) resolves `BIDGO_DATABASE_*` variables via `{{env ...}}` templates, and the `terndotenv` wrapper loads `.env` before invoking tern so those variables exist.

## Architecture

- `cmd/api/` — entrypoint for the HTTP server.
- `cmd/terndotenv/` — migration runner wrapper (loads `.env`, shells out to `tern`).
- `internal/api/` — the `Api` struct (holds the chi router), route bindings in `routes.go`, and HTTP handlers split by domain (`user_handlers.go`, etc.). Routes are nested under `/api/v1/`.
- `internal/store/pgstore/` — database layer. Migrations live in `migrations/`; sqlc reads the migration files as schema plus SQL in `queries/` and generates Go code directly into this package (pgx/v5, `uuid` columns map to `github.com/google/uuid`).

Configuration comes from environment variables prefixed `BIDGO_DATABASE_*`, loaded from `.env` (gitignored; `.env.example` is the template). docker-compose uses the same variables for the Postgres container.
