package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	// "github.com/Resolution-hash/shop_bot/internal/card"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/Resolution-hash/shop_bot/internal/sessions"
	"github.com/fatih/color"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var SessionManager = sessions.NewSessionManager()

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
			color.Blue("\n\n\n" + update.Message.Text)
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

func deleteMessage(bot *tgbotapi.BotAPI, chatInfo *sessions.ChatInfo) {
	deleteConfig := tgbotapi.NewDeleteMessage(chatInfo.ChatID, chatInfo.MessageID)
	if _, err := bot.Send(deleteConfig); err != nil {
		log.Printf("Ошибка при удалении сообщения: %s\n", err)
	}
}

func getKeyboard(value string, back interface{}) tgbotapi.InlineKeyboardMarkup {
	fmt.Println(value)
	switch value {
	case "start":
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

func getMessageText(step string) string {
	switch step {
	case "start":
		return "Приветствие!"
	default:
		return "Такой команды нет. Пожалуйста, выберите из доступных команд"
	}
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userInfo := getUserInfo(update)
	chatInfo := getChatInfo(update)
	deleteMessage(bot, chatInfo)
	switch update.Message.Text {
	case "/start":
		keyboard := getKeyboard("start", nil)
		messageText := getMessageText("start")
		SessionManager.CreateSession(userInfo, chatInfo, keyboard, "start", "start")
		SessionManager.PrintSessionByID(userInfo.UserID)
		sendMessage(bot, update.Message.Chat.ID, messageText, keyboard)
	default:
		keyboard := getKeyboard("start", nil)
		messageText := getMessageText("")
		SessionManager.CreateSession(userInfo, chatInfo, keyboard, "errorCommand", "")
		SessionManager.PrintSessionByID(userInfo.UserID)
		sendMessage(bot, update.Message.Chat.ID, messageText, keyboard)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	userInfo := getUserInfo(update)
	chatInfo := getChatInfo(update)
	_, exists := SessionManager.GetSession(userInfo.UserID)
	if !exists {
		SessionManager.CreateSession(userInfo, chatInfo, tgbotapi.NewInlineKeyboardMarkup(), data, "")
	}
	deleteMessage(bot, chatInfo)

	switch data {
	case "ceramic":
		keyboard := getKeyboard(data, nil)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintSessionByID(userInfo.UserID)
		SessionManager.UpdatePrevStep(userInfo.UserID, data)
		sendMessage(bot, chatInfo.ChatID, "Выберите категорию: ", keyboard)
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

		// sendMessage(bot, chatID, messageText, cardSession.Keyboard)

		fmt.Println("\n\n\n\n", products)
	case "grenades":
		// db, err := setupDatabase()
		// if err != nil {
		// 	log.Println(err)
		// }
		// service := initServices(db)

		// products, err := service.GetProductByType(data)
		// if err != nil {
		// 	log.Println(err)
		// }

		// sessionClose(chatID)

		// cardSession := initCardSession(chatID, products, "ceramic")
		// messageText := cardSession.Card.GetTextTemplate()
		// sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "drawings":
		// db, err := setupDatabase()
		// if err != nil {
		// 	log.Println(err)
		// }
		// service := initServices(db)

		// products, err := service.GetProductByType(data)
		// if err != nil {
		// 	log.Println(err)
		// }
		// sessionClose(chatID)

		// cardSession := initCardSession(chatID, products, "ceramic")
		// messageText := cardSession.Card.GetTextTemplate()
		// sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "showAllItems":
		// db, err := setupDatabase()
		// if err != nil {
		// 	log.Println(err)
		// }
		// service := initServices(db)

		// products, err := service.GetAllProducts()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// //Если уже существует, закрываем
		// sessionClose(chatID)
		// // Создаем сессию для конкретного чата для управления текущим состоянием карточки
		// cardSession := initCardSession(chatID, products, "ceramic")
		// messageText := cardSession.Card.GetTextTemplate()
		// sendMessage(bot, chatID, messageText, cardSession.Keyboard)
	case "prev":

	case "next":

	}

}

func getUserInfo(update tgbotapi.Update) *sessions.UserInfo {
	if update.CallbackQuery != nil {
		user := update.CallbackQuery.From
		return &sessions.UserInfo{
			UserID:     user.ID,
			First_name: user.FirstName,
			Last_name:  user.LastName,
			User_name:  user.UserName,
		}

	}
	user := update.Message.From
	return &sessions.UserInfo{
		UserID:     user.ID,
		First_name: user.FirstName,
		Last_name:  user.LastName,
		User_name:  user.UserName,
	}
}

func getChatInfo(update tgbotapi.Update) *sessions.ChatInfo {
	if update.CallbackQuery != nil {
		return &sessions.ChatInfo{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		}
	}
	return &sessions.ChatInfo{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	}
}
