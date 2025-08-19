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

The service uses **Viper** for configuration management, supporting both TOML files and environment variables with automatic fallbacks:

### Public Configuration (`config.toml`)

Public configuration values that can be safely committed to git:

```toml
[server]
addr = ":8081"
shutdown_timeout_seconds = 30

[jwt]
issuer = "auth.quiby.ai"
audience = "api.quiby.ai"
access_ttl_minutes = 15
```

### Environment Variable Overrides

You can override any config value using environment variables:

| Environment Variable | Config Key | Description |
|---------------------|------------|-------------|
| `SERVER_ADDR` | `server.addr` | Server address and port |
| `SHUTDOWN_TIMEOUT_SECONDS` | `server.shutdown_timeout_seconds` | Graceful shutdown timeout |
| `JWT_ISSUER` | `jwt.issuer` | JWT issuer claim |
| `JWT_AUDIENCE` | `jwt.audience` | JWT audience claim |
| `JWT_ACCESS_TTL` | `jwt.access_ttl_minutes` | JWT access token TTL |

### Private Configuration (Environment Variables)

Sensitive configuration that should be kept private:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `PG_DSN` | PostgreSQL connection string | **Yes** | - |
| `JWT_SECRET_B64` | Base64-encoded JWT secret | **Yes** | - |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | **Yes** | - |

### Optional Overrides

These environment variables can override the defaults in `config.toml`:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | Server address and port | `:8081` |
| `SHUTDOWN_TIMEOUT_SECONDS` | Graceful shutdown timeout | `30` |
| `JWT_ISSUER` | JWT issuer claim | `auth.quiby.ai` |
| `JWT_AUDIENCE` | JWT audience claim | `api.quiby.ai` |
| `JWT_ACCESS_TTL` | JWT access token TTL | `15m` |

### Configuration Precedence

Viper follows this order of precedence (highest to lowest):
1. **Environment Variables** - Override everything
2. **Config File** (`config.toml`) - Default values
3. **Hard-coded defaults** - Fallback values

### Setup Configuration

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your actual values:**
   ```env
   # Database Configuration
   PG_DSN=postgres://user:pass@localhost:5432/auth?sslmode=disable
   
   # JWT Configuration
   JWT_SECRET_B64=your_base64_encoded_secret_here
   
   # Telegram Configuration
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   ```

3. **For production deployments:**
   - Set environment variables directly in your deployment platform
   - Use GitHub Secrets for CI/CD pipelines
   - Never commit `.env` files to version control

### Configuration Helper Functions

The config package provides helper functions for accessing configuration values:

```go
import "github.com/quiby-ai/auth-service/config"

// Get string values
serverAddr := config.GetString("server.addr")

// Get int values  
timeout := config.GetInt("server.shutdown_timeout_seconds")

// Get duration values
ttl := config.GetDuration("jwt.access_ttl_minutes")

// Check if a key is set
if config.IsSet("custom.key") {
    // Handle custom configuration
}
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
│   └── main.go            # Application entry point
├── config/
│   └── config.go          # Configuration management
├── internal/
│   ├── database/
│   │   └── postgres.go    # Database connection
│   ├── handler/
│   │   └── auth.go        # HTTP handlers
│   └── models/
│       └── user.go        # User model and repository
├── db/
│   └── migrations/
│       └── 0001_init.sql  # Database schema
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
├── Makefile               # Build and development commands
└── README.md              # This file
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
