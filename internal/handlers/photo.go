package handlers

import (
	"fmt"
	"io"

	"github.com/Resolution-hash/shop_bot/internal/storage"
	"github.com/Resolution-hash/shop_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

func UploadPhotos(bot *tgbotapi.BotAPI, photo tgbotapi.PhotoSize) (string, error) {
	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		color.Redln("error to get file")
		return "", err
	}

	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	response, err := bot.Client.Get(url)
	if err != nil {
		color.Redln("error downloading file")
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		color.Redln("error to reading file data")
		return "", err
	}

	objectName := utils.GenereteUniqueFileName(".jpg")

	err = storage.MinIOPutPhoto(objectName, data)
	if err != nil {
		return "", err
	}

	return objectName, nil

}

func RemovePhotoFromStorage(filename string) error {
	err := storage.MinIORemovePhoto(filename)
	if err != nil {
		return err
	}
	return nil
}
