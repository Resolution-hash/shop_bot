package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/card"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var sessions map[int64]card.CardSession

func StartBot(cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Println("Bot launched!")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Println("\n\n\n" + update.Message.Text)
			handleCommand(bot, update)
		}
		if update.CallbackQuery != nil {
			handleCallback(bot, update)
		}

		//keyboard := tgbotapi.NewReplyKeyboard(
		//	tgbotapi.NewKeyboardButtonRow(
		//		tgbotapi.NewKeyboardButton("Керамика"),
		//		tgbotapi.NewKeyboardButton("Еще хуйня"),
		//	),
		//)
		//message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//message.ReplyMarkup = keyboard
		//bot.Send(message)
		//if update.CallbackQuery != nil {
		//	handleCallback(update, bot)
		//}
		//if update.Message != nil {
		//	log.Println("\n\n\n" + update.Message.Text)
		//	handleCommand(update, bot)
		//	handleCommand(update, bot)
		//}

	}
}

func initCardSession(chatID int64, products []repository.Product, prevStage string) *card.CardSession {
	if sessions == nil {
		sessions = make(map[int64]card.CardSession)
	}

	currentSession, exists := sessions[chatID]
	if !exists || currentSession.Card == nil {
		currentSession = card.CardSession{
			Card:     card.NewCard(products),
			Keyboard: selectKeyboard("card", prevStage),
		}
		sessions[chatID] = currentSession
	}

	return &currentSession

}

func sessionClose(chatID int64) {
	if _, ok := sessions[chatID]; !ok {
		return
	}
	delete(sessions, chatID)
}
func setupDatabase() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		return nil, err
		//keyboard := selectKeyboard("showAllItems")
		//sendMessage(bot, chatID, err.Error(), keyboard)
	}
	return db, nil
}

func initServices(db *sql.DB) *services.ProductService {
	repo := repository.NewSqliteProductRepo(db)
	service := services.NewProductService(repo)
	return service
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard.(tgbotapi.InlineKeyboardMarkup)
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %s\n", err)
	}

}

func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := bot.Send(deleteConfig); err != nil {
		log.Printf("Ошибка при удалении сообщения: %s\n", err)
	}
}

func selectKeyboard(value string, back interface{}) tgbotapi.InlineKeyboardMarkup {
	fmt.Println(value)
	switch value {
	case "startKeyboard":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Керамика", "ceramic"),
			),
		)
		return keyboard
	case "ceramic":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Гранаты", "grenades"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Лимоны", "lemons"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Рисунки", "drawings"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать все", "showAllItems"),
			),
		)
		return keyboard
	case "card":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("←", "prev"),
				tgbotapi.NewInlineKeyboardButtonData("→", "next"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться", back.(string)),
			),
		)
		return keyboard

	default:
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		log.Println("value is not found on func selectKeyboard()")
		return keyboard
	}
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		keyboard := selectKeyboard("startKeyboard", nil)
		sendMessage(bot, update.Message.Chat.ID, "Приветствие!", keyboard)
	default:
		keyboard := selectKeyboard("startKeyboard", nil)
		sendMessage(bot, update.Message.Chat.ID, "Такой команды нет. Пожалуйста, выберите из доступных товаров!", keyboard)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	deleteMessage(bot, chatID, messageID)

	switch data {
	case "ceramic":

		keyboard := selectKeyboard(data, nil)
		sendMessage(bot, chatID, "Выберите категорию: ", keyboard)
	case "lemons":
		db, err := setupDatabase()

		if err != nil {
			log.Println(err)
		}
		service := initServices(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}
		//Если уже существует, закрываем
		sessionClose(chatID)

		cardSession := initCardSession(chatID, products, "ceramic")
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)

		fmt.Println("\n\n\n\n", products)
	case "grenades":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initServices(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}
		sessionClose(chatID)

		cardSession := initCardSession(chatID, products, "ceramic")
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "drawings":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initServices(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}
		sessionClose(chatID)

		cardSession := initCardSession(chatID, products, "ceramic")
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "showAllItems":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initServices(db)

		products, err := service.GetAllProducts()
		if err != nil {
			log.Fatal(err)
		}

		//Если уже существует, закрываем
		sessionClose(chatID)
		// Создаем сессию для конкретного чата для управления текущим состоянием карточки
		cardSession := initCardSession(chatID, products, "ceramic")
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "prev":
		cardSession, exists := sessions[chatID]
		if !exists {
			return
		}
		cardSession.Card.Prev()
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "next":
		cardSession, exists := sessions[chatID]
		if !exists {
			return
		}
		cardSession.Card.Next()
		messageText := cardSession.Card.GetTextTemplate()
		sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	}

}
