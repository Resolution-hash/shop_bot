package handlers

import (
	"fmt"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

func HandleUploadPhotos(bot *tgbotapi.BotAPI, photo tgbotapi.PhotoSize, caption string) {
	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		color.Redln("error to get file", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	response, err := bot.Client.Get(url)
	if err != nil {
		color.Redln("error downloading file", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil{
		color.Redln("error to reading file data", err)
	}

}
