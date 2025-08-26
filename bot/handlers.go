package bot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hvmidrezv/gozipperbot/utils"
)

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		b.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update.CallbackQuery)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() && message.Command() == "start" {
		b.sendStartMessage(message.Chat.ID)
		return
	}

	if message.Document != nil {
		b.handleDocument(message)
		return
	}
}

func (b *Bot) handleCallbackQuery(cq *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cq.ID, "")
	if _, err := b.API.Request(callback); err != nil {
		log.Printf("Error sending callback response: %v", err)
	}

	chatID := cq.Message.Chat.ID
	state := b.State.Get(chatID)

	switch cq.Data {
	case "compress_action":
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("compress-%d-*", chatID))
		if err != nil {
			log.Printf("Error creating temp dir: %v", err)
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "خطایی در سرور رخ داد."))
			return
		}

		state.State = "awaiting_files"
		state.TempDir = tempDir
		b.State.Set(chatID, state)

		msgText := "حالت فشرده‌سازی فعال شد.\n\n"
		msgText += "لطفا فایل‌های خود را یکی پس از دیگری ارسال کنید.\n"
		msgText += "پس از ارسال تمام فایل‌ها، روی دکمه «پایان و فشرده‌سازی» کلیک کنید."

		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ پایان و فشرده‌سازی", "finish_compress"),
			),
		)
		_, _ = b.API.Send(msg)

	case "decompress_action":
		state.State = "awaiting_zip"
		b.State.Set(chatID, state)
		_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "حالت استخراج فعال شد. لطفا فایل `.zip` خود را ارسال کنید."))

	case "finish_compress":
		if state.State != "awaiting_files" {
			return
		}

		tempDir := state.TempDir
		state.State = "idle"
		b.State.Set(chatID, state)

		files, err := os.ReadDir(tempDir)
		if err != nil || len(files) == 0 {
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "هیچ فایلی برای فشرده‌سازی ارسال نشده است. عملیات لغو شد."))
			_ = os.RemoveAll(tempDir)
			return
		}

		_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "در حال فشرده‌سازی فایل‌ها..."))

		zipFileName := filepath.Join(os.TempDir(), fmt.Sprintf("archive-%d.zip", chatID))
		err = utils.ZipSource(tempDir, zipFileName)
		if err != nil {
			log.Printf("Error zipping source: %v", err)
			return
		}

		doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(zipFileName))
		doc.Caption = "فایل فشرده شما آماده است."
		_, _ = b.API.Send(doc)

		_ = os.Remove(zipFileName)
		_ = os.RemoveAll(tempDir)
	}
}

func (b *Bot) handleDocument(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	state := b.State.Get(chatID)

	switch state.State {
	case "awaiting_files":
		downloadPath := filepath.Join(state.TempDir, message.Document.FileName)
		err := utils.DownloadFile(b.API, message.Document.FileID, downloadPath)
		if err != nil {
			log.Printf("Error downloading file: %v", err)
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "خطایی در دانلود فایل رخ داد."))
			return
		}
		_, _ = b.API.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("فایل `%s` دریافت شد.", message.Document.FileName)))

	case "awaiting_zip":
		if filepath.Ext(message.Document.FileName) != ".zip" {
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "لطفا فقط فایل با فرمت .zip ارسال کنید."))
			return
		}
		_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "فایل zip دریافت شد. در حال استخراج..."))

		tempZipPath := filepath.Join(os.TempDir(), message.Document.FileName)
		err := utils.DownloadFile(b.API, message.Document.FileID, tempZipPath)
		if err != nil {
			log.Printf("Error downloading zip: %v", err)
			return
		}

		destDir := filepath.Join(os.TempDir(), fmt.Sprintf("unzipped-%d", chatID))
		_ = os.MkdirAll(destDir, os.ModePerm)

		err = utils.UnzipSource(tempZipPath, destDir)
		if err != nil {
			log.Printf("Error unzipping: %v", err)
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "خطایی در استخراج فایل رخ داد."))
			return
		}

		files, err := os.ReadDir(destDir)
		if err != nil {
			log.Printf("Error reading extracted files: %v", err)
			_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "خطایی در خواندن فایل‌های استخراج شده رخ داد."))
			return
		}

		for _, file := range files {
			filePath := filepath.Join(destDir, file.Name())
			if !file.IsDir() {
				doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filePath))
				doc.Caption = fmt.Sprintf("فایل: %s", file.Name())
				_, _ = b.API.Send(doc)
			}
		}

		state.State = "idle"
		b.State.Set(chatID, state)
		_ = os.Remove(tempZipPath)
		_ = os.RemoveAll(destDir)

	default:
		b.sendStartMessage(chatID)
	}
}

func (b *Bot) sendStartMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "سلام! به ربات فشرده‌ساز خوش آمدید. لطفا عملیات مورد نظر را انتخاب کنید:")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗜️ فشرده‌سازی", "compress_action"),
			tgbotapi.NewInlineKeyboardButtonData("📂 استخراج", "decompress_action"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, _ = b.API.Send(msg)
}
