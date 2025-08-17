# Auth Service

A minimal Telegram Mini App authentication service built in Go with PostgreSQL.

## Features

- **Telegram Mini App Authentication**: Validates `initData` and issues short-lived JWT tokens
- **User Management**: UPSERTs users with full Telegram profile data
- **JWT-based Authorization**: Uses HS256 signing with configurable TTL
- **Database Persistence**: Stores user data in PostgreSQL with JSONB profile storage
- **Health Monitoring**: Includes health check endpoint with database connectivity testing

## API Endpoints

- `POST /auth/telegram/login` - Authenticate with Telegram initData
- `GET /me` - Get current user profile (protected)
- `POST /auth/logout` - Logout (clears cookie)
- `GET /healthz` - Health check with database ping

## Configuration

The service requires the following environment variables:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SERVER_ADDR` | Server address and port | No | `:8080` |
| `SHUTDOWN_TIMEOUT_SECONDS` | Graceful shutdown timeout | No | `30` |
| `PG_DSN` | PostgreSQL connection string | **Yes** | - |
| `JWT_ISSUER` | JWT issuer claim | No | `auth.quiby.ai` |
| `JWT_AUDIENCE` | JWT audience claim | No | `api.quiby.ai` |
| `JWT_ACCESS_TTL` | JWT access token TTL | No | `15m` |
| `JWT_SECRET_B64` | Base64-encoded JWT secret | **Yes** | - |
| `COOKIE_NAME` | Access cookie name | No | `cp_at` |
| `COOKIE_DOMAIN` | Cookie domain | No | - |
| `COOKIE_PATH` | Cookie path | No | `/` |
| `COOKIE_SECURE` | Secure cookie flag | No | `true` |
| `COOKIE_SAMESITE` | Cookie SameSite policy | No | `none` |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | **Yes** | - |

### Example Configuration

Create a `.env` file in the project root:

```env
# Server Configuration
SERVER_ADDR=:8080
SHUTDOWN_TIMEOUT_SECONDS=30

# Database Configuration
PG_DSN=postgres://user:pass@localhost:5432/auth?sslmode=disable

# JWT Configuration
JWT_ISSUER=auth.quiby.ai
JWT_AUDIENCE=api.quiby.ai
JWT_ACCESS_TTL=15m
JWT_SECRET_B64=your_base64_encoded_secret_here

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

## Database Setup

### 1. Create Database

```sql
CREATE DATABASE auth;
```

### 2. Run Migration

```bash
make migrate-up
```

Or manually:

```bash
psql "$PG_DSN" -f db/migrations/0001_init.sql
```

## Running the Service

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Valid Telegram bot token

### Development

```bash
# Install dependencies
make deps

# Run the service
make run
```

### Production

```bash
# Build binary
make build

# Run binary
./auth-service
```

## Usage

### 1. Telegram Login

Send a POST request to `/` with the `Authorization` header:

```
Authorization: tma <init-data>
```

The service will:
1. Validate the Telegram `initData`
2. UPSERT the user in the database
3. Issue a JWT token
4. Return `user_id` and `ok: true`

### 2. Protected Endpoints

For protected endpoints like `/me`, include the JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt-token>
```

## Development

### Available Make Commands

- `make run` - Run the service
- `make test` - Run tests
- `make lint` - Run linter
- `make build` - Build binary
- `make clean` - Remove binary
- `make deps` - Install dependencies
- `make migrate-up` - Run database migration

### Project Structure

```
auth-service/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   └── postgres.go      # Database connection
│   ├── handler/
│   │   └── auth.go          # HTTP handlers
│   └── models/
│       └── user.go          # User model and repository
├── db/
│   └── migrations/
│       └── 0001_init.sql    # Database schema
├── go.mod                   # Go module file
├── go.sum                   # Go module checksums
├── Makefile                 # Build and development commands
└── README.md                # This file
```

## Security Notes

- JWT tokens are short-lived (default: 15 minutes)
- No refresh tokens - users must re-authenticate
- All database operations use parameterized queries
- Telegram initData validation includes timeout checks
- Bot users are explicitly forbidden

## Dependencies

- `github.com/quiby-ai/common/pkg/auth` - Authentication utilities
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/joho/godotenv` - Environment variable loading
