# SleepTracker Architecture

## Request flow

```
HTTP Request
    → Gin Router
    → Middleware (CORS, logging, recovery, JWT)
    → Handler
    → Service (business logic)
    → Repository
    → PostgreSQL
```

## Package layout

| Package | Purpose |
|---------|---------|
| `cmd/tracker` | Entry point; starts app, handles shutdown |
| `internal/app` | Loads config, connects DB, wires services, starts HTTP server |
| `internal/config` | Environment config via `goloop/env` |
| `internal/transport/handler` | Gin handlers, routes, middleware |
| `internal/service` | Auth (JWT), user, friends, sleep business logic |
| `internal/repo` | Repository interfaces and implementations |
| `models` | Domain types, request/response DTOs, validation helpers |
| `pkg/date` | `date.Date` with `YYYY-MM-DD` JSON support |
| `pkg/score` | Sleep score calculation (age, duration, times) |
| `pkg/tokens` | JWT creation and parsing (HS256) |
| `pkg/logger` | Zap logging (console and file) |

## Database schema

### users

| Column | Type | Notes |
|--------|------|-------|
| id | bigserial | PK |
| username | varchar(40) | unique |
| nickname | text | required |
| password_hash | bytea | bcrypt |
| birth | date | required |
| created_at | timestamptz | default now() |

### friends

| Column | Type | Notes |
|--------|------|-------|
| id | bigserial | PK |
| user1_id | bigint | FK → users |
| user2_id | bigint | FK → users |
| is_accepted | boolean | default false |
| created_at | timestamptz | |
| updated_at | timestamptz | |
| unique (user1_id, user2_id) | | |
| check user1_id ≠ user2_id | | |

### sleeps

| Column | Type | Notes |
|--------|------|-------|
| id | bigserial | PK |
| user_id | bigint | FK → users |
| sleep_time | smallint | minutes 0–2159, nullable |
| wake_time | smallint | minutes 0–2159, nullable |
| score | smallint | 0–100 |
| enter_date | date | unique per user |
| created_at | timestamptz | |
| check sleep_time < wake_time (when both set) | | |

## Sleep time format

Sleep and wake times are stored as **minutes since midnight** (0–2159). This supports times that cross midnight (e.g. sleep 23:00, wake 07:00).

## Sleep score

The score (0–100) is computed in `pkg/score` using:

- User age (from birth)
- Sleep duration
- Bedtime and wake time

## Migrations

SQL migrations live in `migrations/` and are run with [Goose](https://github.com/pressly/goose). In Docker, the `migrations` service runs `goose up` before the app starts.
