# Telegram AML Bot

A Telegram bot for checking cryptocurrency addresses and transactions against AML (Anti-Money Laundering) databases.

## Features

- Check cryptocurrency addresses for suspicious activity
- Check transaction hashes for AML compliance
- Real-time results with risk scores
- Detailed reporting of suspicious activities
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Telegram Bot Token (from [@BotFather](https://t.me/BotFather))
- AML Provider API Key

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/clevertechru/tgbot_aml.git
cd tgbot_aml
```

2. Create and configure `.env` file:
```bash
cp .env.example .env
# Edit .env with your actual API keys
```

3. Build and run with Docker:
```bash
docker compose up --build
```

## Configuration

### Environment Variables

Create a `.env` file with the following variables:

```env
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here
AML_API_KEY=your_api_key_here

# Optional
AML_BASE_URL=https://api.aml-provider.com
```

### Bot Commands

- `/start` - Start the bot and get welcome message
- `/check <address>` - Check a cryptocurrency address
- `/check <tx_hash>` - Check a transaction hash

## Development

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run the bot:
```bash
go run cmd/bot/main.go
```

### Building

```bash
# Build binary
go build -o bot ./cmd/bot

# Build Docker image
docker build -t tgbot_aml .
```

### Testing

```bash
go test ./...
```

## Docker Deployment

### Production

1. Update `.env` with production credentials
2. Run in detached mode:
```bash
docker compose up -d
```

### Development

Run with logs:
```bash
docker compose up --build
```

### Logs

View logs:
```bash
docker compose logs -f
```

## Project Structure

```
.
├── cmd/
│   └── bot/           # Main application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── domain/        # Core domain models and interfaces
│   ├── handlers/      # Telegram bot handlers
│   └── services/      # Business logic services
├── config/            # Configuration files
├── logs/             # Application logs
├── Dockerfile        # Docker build configuration
├── docker-compose.yml # Docker Compose configuration
└── .env              # Environment variables
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 