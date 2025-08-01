# auth-service
The [auth-service] is a standalone microservice responsible for handling all authentication-related functionality across the platform.

## Configuration

The service requires the following environment variables:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `BOT_TOKEN` | Telegram Bot Token (get from @BotFather) | Yes | - |
| `JWT_SECRET` | Secret key for signing JWT tokens | Yes | - |
| `PORT` | Port where the service will run | No | 8081 |
| `SHUTDOWN_TIMEOUT_SECONDS` | Graceful shutdown timeout in seconds | No | 30 |

### Example Configuration

You can set environment variables either via export commands or by creating a `.env` file:

**Option 1: Environment Variables**
```bash
export BOT_TOKEN="your_telegram_bot_token_here"
export JWT_SECRET="your_jwt_secret_here"
export PORT="8081"
export SHUTDOWN_TIMEOUT_SECONDS="30"
```

**Option 2: .env File**
Create a `.env` file in the project root:
```env
BOT_TOKEN="your_telegram_bot_token_here"
JWT_SECRET="your_jwt_secret_here"
PORT="8081"
SHUTDOWN_TIMEOUT_SECONDS="30"
```

## Running the Service

```bash
go run cmd/main.go
```
