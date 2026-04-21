# Agritech Monitoring Webhook Server

A specialized monitoring middleware built with Go, designed to bridge Grafana alerts and automated notification systems for smart farm infrastructure.

## Key Features
- **Real-time Alert Processing**: Handles incoming HTTP POST webhooks from Grafana.
- **Dynamic Messaging**: Parses alert payloads to deliver context-aware notifications.
- **Automated Response**: Triggers cooling system logic based on farm temperature thresholds.
- **Secure Configuration**: Uses environment variables (.env) for credential management.

## Tech Stack
- **Language**: Go 1.20+
- **Monitoring**: Grafana
- **Alerting**: Telegram Bot API
- **Infrastructure**: Linux systemd (Daemonized)

## Architecture Flow
1. Farm sensors send data to InfluxDB.
2. Grafana monitors threshold (e.g., Temp > 30C).
3. Grafana triggers Webhook to this Go Server.
4. Go Server processes logic and notifies via Telegram.

## Setup
1. Configure `.env` with your `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`.
2. Build: `go build -o webhook-server main.go`
3. Run as a service using the provided systemd configuration.
