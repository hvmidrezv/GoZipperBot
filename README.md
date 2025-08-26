# 🤖 GoZipperBot - Telegram File Archiver
![img.png](img.png)
<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8.svg?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Docker-20.10+-2496ED.svg?style=for-the-badge&logo=docker" alt="Docker Version">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge" alt="License: MIT">
</p>

GoZipperBot is a fully-functional Telegram bot written in **Go**, allowing users to compress and extract `.zip` files directly within Telegram. It features a smooth and interactive experience with inline keyboards and supports password-protected archives.

This project was created as a **portfolio piece** to showcase backend development skills in Go, stateful bot design, API integration, and Docker-based deployment.

---

## ✨ Features

- ✔️ **File Compression**  
  Send multiple photos or documents and receive a single `.zip` archive in return.

- ✔️ **File Extraction**  
  Send a `.zip` file to the bot, and it will extract and return its contents.

- ⚙️ **Interactive Workflow**
    - All operations use inline keyboards.
    - Shows a list of collected files before finalizing.
    - Users can cancel an operation at any time.

- 🚀 **Dockerized**  
  Comes with a multi-stage Dockerfile for lightweight, production-ready builds.

---

## 🛠️ Tech Stack

- **Language**: Go
- **Container Platform**: Docker

**Key Go Libraries:**

- [`go-telegram-bot-api`](https://github.com/go-telegram-bot-api/telegram-bot-api) — Telegram Bot API integration
- [`alexmullins/zip`](https://github.com/alexmullins/zip) — Password-protected zip handling
- Go Standard Libraries: `archive/zip`, `os`, `path/filepath`

---

## 🚀 Setup and Usage

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/zipperbot.git
cd zipperbot
```

### 2. Get a Telegram Bot Token
Talk to [@BotFather](https://t.me/BotFather) and generate a new token.

---

### 🧪 Method 1: Run Locally with Go

#### On Linux / macOS:

```bash
export TELEGRAM_BOT_TOKEN="YOUR_TOKEN_HERE"
go run main.go
```

#### On Windows (CMD):

```bash
set TELEGRAM_BOT_TOKEN="YOUR_TOKEN_HERE"
go run main.go
```

---

### 🐳 Method 2: Run with Docker

Build the image:

```bash
docker build -t telegram-zipper-bot .
```

Run the bot:

```bash
docker run --rm -it -e TELEGRAM_BOT_TOKEN="YOUR_TOKEN_HERE" telegram-zipper-bot
```

---

## 🤖 How to Use

1. Start the bot on Telegram (`/start`)
2. Choose "Compress" or "Decompress" using inline buttons.
3. Follow the instructions to complete the process.

---

## 📁 Project Structure

```
/
├── go.mod
├── main.go                # Entry point
├── Dockerfile             # Docker build instructions
│
├── bot/                   # Bot logic
│   ├── bot.go             # Bot setup & update loop
│   ├── handlers.go        # Message & file handling
│   └── state.go           # Per-user state manager
│
└── utils/                 # Utilities
    ├── archive.go         # Compression / extraction logic
    └── downloader.go      # Telegram file downloads
```

---

## 📄 License

Licensed under the [MIT License](LICENSE).