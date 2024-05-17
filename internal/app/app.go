package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"


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

func sendCard(bot *tgbotapi.BotAPI, userID int, text string, keyboard tgbotapi.InlineKeyboardMarkup, imageName string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	path := cfg.ImagesUrl + "\\" + imageName + ".jpg"
	color.Redln(path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error to upload file")
	}
	defer file.Close()

	card := tgbotapi.NewPhotoUpload(int64(userID), file.Name())
	card.File = tgbotapi.FileReader{
		Name:   file.Name(),
		Reader: file,
		Size:   -1,
	}
	card.Caption = text
	card.ReplyMarkup = keyboard
	bot.Send(card)
}

func deleteMessage(bot *tgbotapi.BotAPI, messageID int, userID int) {
	deleteConfig := tgbotapi.NewDeleteMessage(int64(userID), messageID)
	if _, err := bot.Send(deleteConfig); err != nil {
		log.Printf("Ошибка при удалении сообщения: %s\n", err)
	}
}

func getKeyboard(value string, back interface{}) tgbotapi.InlineKeyboardMarkup {
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
				tgbotapi.NewInlineKeyboardButtonData("Посуда для питья", "drinkware"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Посуда для еды", "dishware"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать все", "showAllItems"),
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

func getDynamicKeyboard(value string, prevStep string, quantity string) tgbotapi.InlineKeyboardMarkup {
	switch value {
	case "addToCart":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("←", "prev"),
				tgbotapi.NewInlineKeyboardButtonData("→", "next"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("-", "increment"),
				tgbotapi.NewInlineKeyboardButtonData(quantity, "no_action"),
				tgbotapi.NewInlineKeyboardButtonData("+", "decrement"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться", prevStep),
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
				tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", "addToCart"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться", prevStep),
			),
		)
		return keyboard
	case "back":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", prevStep),
			),
		)
		return keyboard
	default:
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		color.Redln("value is not found on func getKeyboard()")
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
	case "drinkware":
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
		keyboard := getDynamicKeyboard("card", session.PrevStep, "")
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendCard(bot, userInfo.UserID, card.GetTextTemplate(), keyboard, card.Image)
	case "dishware":
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
		keyboard := getDynamicKeyboard("card", session.PrevStep, "")
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendCard(bot, userInfo.UserID, card.GetTextTemplate(), keyboard, card.Image)
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
		sendCard(bot, userInfo.UserID, card.GetTextTemplate(), keyboard, card.Image)
	case "showAllItems":
		db, err := setupDatabase()
		if err != nil {
			color.Redln(err)
		}
		service := initProductService(db)

		products, err := service.GetAllProducts()
		if err != nil {
			color.Redln(err)
		}

		if repository.IsEmpty(products) {
			keyboard := getKeyboard("back", session.PrevStep)
			SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
			sendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)
		keyboard := getDynamicKeyboard("card", session.PrevStep, "")
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		SessionManager.PrintLogs(userInfo.UserID)
		sendCard(bot, userInfo.UserID, card.GetTextTemplate(), keyboard, card.Image)
	case "prev":
		currentCard := session.CardManager.CurrentCard
		defer func() {
			SessionManager.PrintLogs(userInfo.UserID)
			session.CardManager.PrintLogs()
		}()
		color.Blueln("prev")
		currentCard.Prev()
		keyboard := getDynamicKeyboard("card", session.PrevStep, "")
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)

		sendCard(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard, currentCard.Image)
	case "next":
		currentCard := session.CardManager.CurrentCard
		defer func() {
			SessionManager.PrintLogs(userInfo.UserID)
			session.CardManager.PrintLogs()
		}()
		color.Blueln("next")
		currentCard.Next()
		keyboard := getDynamicKeyboard("card", session.PrevStep, "")
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		sendCard(bot, userInfo.UserID, session.CardManager.CurrentCard.GetTextTemplate(), keyboard, currentCard.Image)
	case "addToCart":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)
		currentCard := session.CardManager.CurrentCard

		db, err := setupDatabase()
		if err != nil {
			color.Redln(err)
		}
		defer db.Close()
		repo := repository.NewSqliteCartRepo(db)
		service := services.NewCartService(repo)

		item := repository.CartItem{
			ProductID: currentCard.ID,
			UserID:    int64(userInfo.UserID),
			Quantity:  1,
		}

		_, err = service.AddItem(item)
		if err != nil {
			color.Redln("Error to add item:", err)
		}

		session.CartManager.Add(item)
		keyboard := getDynamicKeyboard("addToCart", session.PrevStep, session.CartManager.Quantitiy(item))
		SessionManager.UpdateSession(userInfo.UserID, keyboard, data)
		sendCard(bot, userInfo.UserID, currentCard.GetTextTemplate(), keyboard, currentCard.Image)
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
