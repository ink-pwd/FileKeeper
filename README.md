# 📦 FileKeeper
FileKeeper — lightweight Go microservice for external file persistence.

It acts as a stateless file relay: files are uploaded to the service, immediately forwarded to Telegram for storage, and never persisted on the server itself. The service returns a Telegram `file_id`, which can later be used to retrieve the file directly from Telegram.

This approach allows offloading file storage outside of the application infrastructure while keeping the system simple, fast, and disposable.

Administration is possible directly from the private Telegram chat, into which we upload files.

---

## 🚀 Features

- Upload files via API
- Automatic forwarding to Telegram bot
- Lightweight microservice architecture
- No local file persistence 
- Dockerized deployment
- Designed for personal / small-server usage

---

## 🧱 Project Structure
```text
FileKeeper/
│
├── cmd/
│ └── api-server/ # Application entry point
│
├── internal/
│ ├── config/ # Config loader
│ ├── consts/ # Basic constants
│ ├── handler/ # Business logic
│ └── telegram/ # Telegram bot integration
│
├── logger/ # Implementing logging
│
├── env # Environment/config files
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## ⚙️ Configuration

FileKeeper is configured via a `.env` file. Copy the example template and fill in your values:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `HOST` | `0.0.0.0` | Network interface to bind to. Use `127.0.0.1` for local-only access. |
| `PORT` | `:8080` | Port the API server listens on. Colon is required. |
| `READTIMEOUT` | `60` | Timeout in seconds for reading the entire incoming request. |
| `WRITETIMEOUT` | `30` | Timeout in seconds for writing the complete response. |
| `MAXFILESIZEMEGABYTE` | `32` | Maximum upload file size in MB. Telegram limits bots to 50 MB. |
| `MAXRAMSIZEMEGABYTE` | `32` | RAM buffer limit in MB for processing files before streaming to Telegram. |
| `TOKEN` | *required* | Telegram Bot API token from [@BotFather](https://t.me/BotFather). **Keep this secret.** |
| `CHATID` | *required* | Your private Telegram chat ID with the bot. See [How to get chat ID](#how-to-get-telegram-credentials). |

### How to get Telegram credentials

1. **Create a bot.** Open [@BotFather](https://t.me/BotFather), send `/newbot`, follow the instructions. Copy the token — this is your `TOKEN`.
2. **Get your chat ID.** Send any message to your new bot, then open this URL in a browser:
   ```
   https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates
   ```
   Find `"chat":{"id":<YOUR_CHAT_ID>}` in the response — this is your `CHATID`.
3. **Fill `.env`** with these values.

### Sample `.env` structure

```env
HOST=0.0.0.0
PORT=:8080
READTIMEOUT=60
WRITETIMEOUT=30
MAXFILESIZEMEGABYTE=32
MAXRAMSIZEMEGABYTE=32

TOKEN=your_bot_token_here
CHATID=your_chat_id_here
```

## 🐳 Run with Docker

### Pull image from Docker Hub

```bash
docker pull inkpwd/filekeeper:latest
```
### Run container
Create `.env` file and `docker-compose.yml` in project directory
RUN
```bash
docker compose up -d
```

## 🔌API Endpoints

### Upload file
the maximum possible file size is configured in .env
```http
POST /upload
```
Form-data:

- file: binary file

Request:

key => file
value => {your_file}.{type}

Response:

```json
{
    "id": "fileID"
}
```
### Get file
```http
GET /files/{id}
```
Response:

```http
file...
200 OK
```
## 🧠 Architecture Idea

FileKeeper acts as a bridge:
```
Client → FileKeeper API → Telegram Bot → Your private storage chat
```

No database required — Telegram acts as storage layer.

## 🧩 Use Case Diagram

(will be added soon)

## 🎥 Video Overview

(will be added soon)

## License
MIT
