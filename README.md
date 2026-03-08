# SleepTracker

A REST API for tracking sleep, managing user profiles, and viewing friends' sleep data.

## Features

- **User management** — Registration, login, profile updates, password changes
- **Sleep tracking** — Record sleep/wake times and computed sleep scores
- **Friends** — Send, accept, and manage friendship requests
- **Friend sleep data** — View friends' sleep history and scores (for accepted friendships)

Sleep scores are computed based on user age, sleep duration, bedtime, and wake time.

## Tech Stack

- **Go 1.24**
- **Gin** — HTTP router
- **PostgreSQL** — Database
- **JWT** — Access and refresh tokens
- **Zap** — Structured logging
- **bcrypt** — Password hashing

## Quick Start

### Local development

```bash
# Install dependencies
go mod tidy

# Set environment variables (required for auth)
export REFRESH_SECRET=your-refresh-secret
export ACCESS_SECRET=your-access-secret

# Optional: override DB connection (default: postgresql://user:password@localhost:5432/app?sslmode=false)
export CONN_STRING=postgresql://user:pass@localhost:5432/tracker?sslmode=disable

# Run
go run ./cmd/tracker/main.go
```

### Docker Compose

```bash
# Create .env with required secrets:
# REFRESH_SECRET=...
# ACCESS_SECRET=...
# Optional: POSTGRES_USER, POSTGRES_PASSWORD, PORT, LOG_PATH

docker-compose up
```

This starts:

1. PostgreSQL on port 5432
2. Goose migrations
3. App on `PORT` (default 8080)

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `CONN_STRING` | `postgresql://user:password@localhost:5432/app?sslmode=false` | PostgreSQL connection string |
| `REFRESH_SECRET` | — | JWT refresh token secret (required) |
| `REFRESH_EXP` | `720h` (Docker) / `24h` (code default) | Refresh token expiry |
| `ACCESS_SECRET` | — | JWT access token secret (required) |
| `ACCESS_EXP` | `24h` | Access token expiry |
| `LOG_PATH` | `./logs` | Log file directory (Docker) |

## API Overview

| Auth | Method | Path | Description |
|------|--------|------|-------------|
| — | GET | `/ping` | Health check |
| — | POST | `/register` | Register user |
| — | POST | `/login` | Login |
| — | POST | `/refresh` | Refresh tokens |
| ✓ | GET | `/users/` | Current user profile |
| ✓ | GET | `/users/:id` | User by ID |
| ✓ | PUT | `/users/` | Update profile |
| ✓ | PATCH | `/users/password` | Change password |
| ✓ | DELETE | `/users/` | Delete account |
| ✓ | POST | `/friends/request` | Send friend request |
| ✓ | POST | `/friends/accept` | Accept friend request |
| ✓ | GET | `/friends/` | List friendships |
| ✓ | DELETE | `/friends/:id` | Delete friendship |
| ✓ | GET | `/friends/:id/history/:page` | Friend's sleep history |
| ✓ | GET | `/friends/:id/table/:page` | Friend's scores |
| ✓ | POST | `/sleep/` | Add sleep record |
| ✓ | PUT | `/sleep/:id` | Update sleep record |
| ✓ | DELETE | `/sleep/:id` | Delete sleep record |
| ✓ | GET | `/history/:page` | Own sleep history |
| ✓ | GET | `/table/:page` | Own scores |

Protected routes require `Authorization: Bearer <access_token>`.

See [docs/API.md](docs/API.md) for detailed request/response formats.

## Project Structure

```
SleepTracker/
├── cmd/tracker/       # Entry point
├── internal/          # Application code
│   ├── app/           # Bootstrap
│   ├── config/        # Environment config
│   ├── repo/          # Data access
│   ├── service/       # Business logic
│   └── transport/     # HTTP handlers
├── migrations/        # Goose SQL migrations
├── models/            # Domain and transport types
├── pkg/               # Reusable packages
│   ├── date/          # Date handling
│   ├── logger/        # Zap logging
│   ├── score/         # Sleep score calculation
│   └── tokens/        # JWT utilities
└── docs/              # Documentation
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for more details.

## License

MIT
