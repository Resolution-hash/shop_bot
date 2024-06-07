package app

import (
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/handlers"
	"github.com/gookit/color"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func StartBot(cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Println("Bot launched!")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			color.Blueln("\n\n\n" + update.Message.Text)
			handlers.HandleCommand(bot, update)
		}
		if update.CallbackQuery != nil {
			color.Blueln("\n\n\n" + update.CallbackQuery.Data)
			handlers.HandleCallback(bot, update)
		}
		if len(*update.Message.Photo) > 1 {
			color.Blueln("\n\n\n" + "upload photo")
			photos := *update.Message.Photo
			photoSize := photos[len(photos)-1]
			handlers.HandleUploadPhotos(bot, photoSize, update.Message.Caption)
		}

	}
}
