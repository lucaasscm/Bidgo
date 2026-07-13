# Bidgo

An auction/bidding HTTP API written in Go.

This is a **study project**. The goal is not the product itself, but learning how to build a real-world Go backend from scratch:

- **HTTP servers** — routing, middleware, handlers, and JSON request/response handling with [chi](https://github.com/go-chi/chi)
- **WebSockets** — real-time communication for live auction bidding
- **Authentication** — password hashing, sessions, and protected routes
- **Go workflow and tools** — project layout, code generation with [sqlc](https://github.com/sqlc-dev/sqlc), migrations with [tern](https://github.com/jackc/tern), hot reload with [air](https://github.com/air-verse/air)

## Stack

| Purpose            | Tool                                  |
| ------------------ | ------------------------------------- |
| Language           | Go 1.26                               |
| Router             | chi v5                                |
| Database           | PostgreSQL 16 (Docker)                |
| DB driver          | pgx v5                                |
| Query codegen      | sqlc                                  |
| Migrations         | tern                                  |
| Hot reload         | air                                   |

## Project layout

```
cmd/
  api/          # HTTP server entrypoint
  terndotenv/   # migration runner (loads .env, invokes tern)
internal/
  api/          # Api struct, routes, HTTP handlers
  store/
    pgstore/    # migrations, SQL queries, sqlc-generated code
```

## Roadmap

- [x] Project scaffold (chi router, docker compose, tern + sqlc + air setup)
- [ ] Users table migration and queries
- [ ] Request validation
- [ ] User signup
- [ ] Login / logout with sessions
- [ ] CSRF protection
- [ ] Auction products CRUD
- [ ] Real-time bidding over WebSockets
