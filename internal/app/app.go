package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/card"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/Resolution-hash/shop_bot/internal/sessions"
	"github.com/gookit/color"

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
			color.Blue.Print("\n\n\n" + update.Message.Text)
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
		fmt.Println("error to get cfg.DbUrl")
		return nil, err
	}
	return db, nil
}

func initProductService(db *sql.DB) *services.ProductService {
	repo := repository.NewSqliteProductRepo(db)
	service := services.NewProductService(repo)
	return service
}

func sendMessage(bot *tgbotapi.BotAPI, userID int, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(int64(userID), text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard.(tgbotapi.InlineKeyboardMarkup)
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %s\n", err)
	}

}

func deleteMessage(bot *tgbotapi.BotAPI, messageID int, userID int) {
	deleteConfig := tgbotapi.NewDeleteMessage(int64(userID), messageID)
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
				tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", "prev"),
				tgbotapi.NewInlineKeyboardButtonData("→", "next"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться", back.(string)),
			),
		)
		return keyboard
	case "back":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", back.(string)),
			),
		)
		return keyboard
	default:
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		log.Println("value is not found on func getKeyboard()")
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
	deleteMessage(bot, userInfo.MessageID, userInfo.UserID)
	switch update.Message.Text {
	case "/start":
		keyboard := getKeyboard("start", nil)
		messageText := getMessageText("start")
		SessionManager.CreateSession(userInfo, keyboard, "start", "start")
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, messageText, keyboard)
	default:
		keyboard := getKeyboard("start", nil)
		messageText := getMessageText("")
		SessionManager.CreateSession(userInfo, keyboard, "errorCommand", "")
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, messageText, keyboard)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	userInfo := getUserInfo(update)
	session, exists := SessionManager.GetSession(userInfo.UserID)
	if !exists {
		SessionManager.CreateSession(userInfo, tgbotapi.NewInlineKeyboardMarkup(), data, "")
	}
	deleteMessage(bot, userInfo.MessageID, userInfo.UserID)

	switch data {
	case "ceramic":
		keyboard := getKeyboard(data, nil)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, "Выберите категорию: ", keyboard)
	case "lemons":

		db, err := setupDatabase()

		if err != nil {
			log.Println(err)
		}
		service := initProductService(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := getKeyboard("back", session.PrevStep)
			SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
			sendMessage(bot, userInfo.UserID, "Товаров нет.", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)

	case "grenades":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initProductService(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := getKeyboard("back", session.PrevStep)
			SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
			sendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)

	case "drawings":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initProductService(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := getKeyboard("back", session.PrevStep)
			SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
			sendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)
	case "showAllItems":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		service := initProductService(db)

		products, err := service.GetAllProducts()
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := getKeyboard("back", session.PrevStep)
			SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
			sendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)
	case "prev":
		defer func() {
			SessionManager.PrintLogs(userInfo.UserID)
			session.CardManager.PrintLogs()
		}()
		color.Blueln("prev")
		session.CardManager.CurrentCard.Prev()
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)

		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)
	case "next":
		defer func() {
			SessionManager.PrintLogs(userInfo.UserID)
			session.CardManager.PrintLogs()
		}()
		color.Blueln("next")
		session.CardManager.CurrentCard.Next()
		keyboard := getKeyboard("card", session.PrevStep)
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		sendMessage(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard)
	}

}

func getUserInfo(update tgbotapi.Update) *sessions.UserInfo {
	if update.CallbackQuery != nil {
		user := update.CallbackQuery.From
		return &sessions.UserInfo{
			UserID:     user.ID,
			MessageID:  update.CallbackQuery.Message.MessageID,
			First_name: user.FirstName,
			Last_name:  user.LastName,
			User_name:  user.UserName,
		}

	}

	user := update.Message.From
	return &sessions.UserInfo{
		UserID:     user.ID,
		MessageID:  update.Message.MessageID,
		First_name: user.FirstName,
		Last_name:  user.LastName,
		User_name:  user.UserName,
	}
}

// func getChatInfo(update tgbotapi.Update) *sessions.ChatInfo {
// 	if update.CallbackQuery != nil {
// 		return &sessions.ChatInfo{
// 			ChatID:    update.CallbackQuery.Message.Chat.ID,
// 			MessageID: update.CallbackQuery.Message.MessageID,
// 		}
// 	}
// 	return &sessions.ChatInfo{
// 		ChatID:    update.Message.Chat.ID,
// 		MessageID: update.Message.MessageID,
// 	}
// }
