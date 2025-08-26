package utils

import (
	"io"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DownloadFile(api *tgbotapi.BotAPI, fileID, destPath string) error {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := api.GetFile(fileConfig)
	if err != nil {
		return err
	}

	downloadURL := file.Link(api.Token)
	response, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
}
