# TestFlight Monitor

<p align="center">
  <img src="https://developer.apple.com/assets/elements/icons/testflight/testflight-96x96_2x.png" width="96" height="96" alt="TestFlight">
</p>

<p align="center">
  <strong>Automatically monitor TestFlight beta slots and get notified instantly</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#docker-deployment">Docker</a> •
  <a href="#configuration">Configuration</a> •
  <a href="README.md">中文</a>
</p>

---

## Features

- 🔍 **Auto Detection** - Periodically check TestFlight availability
- 📱 **Telegram Notifications** - Get notified immediately when slots open
- ⏱️ **Custom Interval** - Adjustable check interval (minimum 10 seconds)
- ♾️ **Forever Monitoring** - Support permanent monitoring without expiration
- 🌐 **Proxy Support** - HTTP/SOCKS5 proxy for restricted networks
- 🌍 **Bilingual UI** - English and Simplified Chinese interface
- 🐳 **Docker Ready** - One-click deployment

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-username/tf-monitor.git
cd tf-monitor

# Start the service
docker-compose up -d

# Access http://localhost:8080
```

### Option 2: Build from Source

```bash
# Build frontend
cd web && npm install && npm run build && cd ..

# Build backend
go build -o tf-monitor ./cmd/server/main.go

# Run
./tf-monitor
```

## Docker Deployment

### docker-compose.yml

```yaml
version: "3.8"

services:
  tf-monitor:
    image: simongino/tf-monitor:latest
    container_name: tf-monitor
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - PROXY_ENABLED=false
      - PROXY_URL=
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | Server port |
| `DB_PATH` | data/tf-monitor.db | Database path |
| `PROXY_ENABLED` | false | Enable proxy |
| `PROXY_URL` | - | Proxy URL, e.g., `http://127.0.0.1:7890` |

## Configuration

### Telegram Notification Setup

1. Send `/newbot` to [@BotFather](https://t.me/BotFather) to create a bot
2. Get the Bot Token (format: `123456789:ABCdefGHI...`)
3. Send a message to [@userinfobot](https://t.me/userinfobot) to get your Chat ID
4. Enter Bot Token and Chat ID in Settings
5. Click "Test Send" to verify

### Proxy Configuration

If you need a proxy to access TestFlight:

```bash
# HTTP proxy
PROXY_URL=http://127.0.0.1:7890

# SOCKS5 proxy
PROXY_URL=socks5://127.0.0.1:7890
```

## Usage

### Adding Monitors

1. Enter TestFlight URLs in the sidebar (one per line for batch)
2. Set check interval (recommended: 60 seconds)
3. Choose monitoring duration (2h/8h/12h/24h/Forever)
4. Click "Add Monitor"

### Notification Modes

| Mode | Description |
|------|-------------|
| Loop | Notify every time slots are available |
| Once | Notify once then stop |
| On Change | Notify only when status changes from "Full" to "Available" |

### Card Actions

- **Pause/Resume** - Pause or resume monitoring
- **Edit** - Modify check interval and duration
- **Delete** - Remove the monitor

## API Reference

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/monitors | List all monitors |
| POST | /api/monitors | Create monitor(s) |
| PUT | /api/monitors/:id | Update monitor |
| DELETE | /api/monitors/:id | Delete monitor |
| POST | /api/monitors/:id/toggle | Toggle monitor |
| GET | /api/telegram | Get Telegram config |
| PUT | /api/telegram | Update Telegram config |
| POST | /api/telegram/test | Test Telegram notification |
| GET | /api/status | Get service status |

## Tech Stack

- **Backend**: Go + Gin + GORM + SQLite
- **Frontend**: Vue 3 + TypeScript + Vite
- **Deployment**: Docker

## License

MIT License

## Acknowledgments

- [Apple TestFlight](https://developer.apple.com/testflight/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
