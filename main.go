package main

import (
	"github.com/hvmidrezv/gozipperbot/bot"
	"log"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	b, err := bot.NewBot(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	b.Start()
}
