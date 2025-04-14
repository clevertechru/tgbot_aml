# Telegram AML Checker Bot

A Telegram bot for checking addresses and transactions against AML (Anti-Money Laundering) databases.

## Features

- Check addresses for suspicious activity
- Check transactions for suspicious activity
- Real-time AML risk assessment
- User-friendly Telegram interface

## Commands

- `/start` - Show welcome message and available commands
- `/check <address>` - Check an address for suspicious activity
- `/checktx <tx_hash>` - Check a transaction for suspicious activity

## Architecture

The project follows a clean architecture with the following components:

- `domain/` - Core business logic and interfaces
  - `aml.go` - AML domain models and interfaces
  - `aml_provider.go` - AML service interface and mock implementation
- `service/` - Business logic implementation
  - `aml.go` - AML service implementation
- `handler/` - HTTP/Telegram handlers
  - `telegram.go` - Telegram bot message handling
- `cmd/bot/` - Application entry point

## Setup

1. Clone the repository
2. Set up environment variables:
   ```bash
   export TELEGRAM_BOT_TOKEN=your_bot_token
   export AML_API_KEY=your_aml_api_key  # Optional, only needed for real AML provider
   ```
3. Run the bot:
   ```bash
   go run cmd/bot/main.go
   ```

## Development

The project uses a mock AML service by default. To implement a real AML provider:

1. Create a new provider implementing the `domain.AMLService` interface
2. Update the service initialization in `cmd/bot/main.go`

## Testing

Run tests with:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License 