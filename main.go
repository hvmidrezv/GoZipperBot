package main

import (
	"github.com/hvmidrezv/gozipperbot/bot"
	"log"
	"net/http"
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

	// Start the bot in a separate goroutine
	go b.Start()

	// Start a dummy HTTP server to bind to the required port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
