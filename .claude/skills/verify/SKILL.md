---
name: verify
description: How to launch and drive the Bidgo API end-to-end for verification.
---

# Verifying Bidgo changes

## Launch

- `docker compose up -d` must already be running (Postgres). `.env` must exist (values are wrapped in double quotes — strip them if parsing in shell).
- Apply migrations with `go run ./cmd/terndotenv` (never `tern migrate` directly).
- Run the server with `go run ./cmd/api` (background). Listens on **:3080**. No output on success besides "Server running on port :3080".

## Drive

All routes are under `http://localhost:3080/api/v1`. Every mutating request needs **both** a session cookie jar and a CSRF token:

```bash
# fresh token per request (also seeds the csrf cookie into the jar)
TOK=$(curl -s -c jar.txt -b jar.txt $BASE/csrftoken | sed 's/.*"csrf_token":"\([^"]*\)".*/\1/')
curl -s -b jar.txt -c jar.txt -X POST -H "X-CSRF-Token: $TOK" -H 'Content-Type: application/json' -d '{...}' $BASE/...
```

Gotchas:
- Signup requires `bio` between 10 and 255 chars, `password` >= 8.
- Product creation requires `auction_end` at least 2h in the future (RFC3339, e.g. `date -u -d '+3 hours' +%Y-%m-%dT%H:%M:%SZ`).
- Auth is session-based (`/users/login` after signup); auth-gated routes return 401 `must be logged in` otherwise.

## Clean up

Delete test users by email — products/bids cascade:

```bash
docker compose exec -T db psql -U <user> -d <db> -c "DELETE FROM users WHERE email IN ('...');"
```
