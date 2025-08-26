package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot ساختار اصلی ربات است که شامل API و مدیر وضعیت می‌باشد
type Bot struct {
	API   *tgbotapi.BotAPI
	State *StateManager
}

// NewBot یک نمونه جدید از ربات را ایجاد و برمی‌گرداند
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		API:   api,
		State: NewStateManager(),
	}, nil
}

// Start ربات را راه‌اندازی کرده و حلقه دریافت آپدیت‌ها را شروع می‌کند
func (b *Bot) Start() {
	b.API.Debug = false
	log.Printf("Authorized on account %s", b.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}
