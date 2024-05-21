package handlers

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/card"
	"github.com/Resolution-hash/shop_bot/internal/messages"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/Resolution-hash/shop_bot/internal/sessions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

var SessionManager = sessions.NewSessionManager()

func HandleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.Message.Text
	userMessageID := update.Message.MessageID
	userInfo := getUserInfo(update)

	session, exists := SessionManager.GetSession(userInfo.UserID)
	if !exists {
		session = SessionManager.CreateSession(userInfo)
	}
	session.LastUserMessageID = userMessageID

	if session.LastUserMessageID != 0 {
		messages.DeleteMessage(bot, session.LastUserMessageID, userInfo.UserID)
	}
	if session.LastBotMessageID != 0 {
		messages.DeleteMessage(bot, session.LastBotMessageID, userInfo.UserID)
	}

	var keyboard tgbotapi.ReplyKeyboardMarkup
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var botMessageID int
	var messageText string
	switch data {
	case "/start":
		keyboard = messages.GetReplyKeyboard()
		messageText = messages.GetMessageText("start")
		messages.SendReplyKeyboard(bot, userInfo.UserID, messageText, keyboard)
	case "Корзина":

		db, err := setupDatabase()
		if err != nil {
			color.Redln(err)
		}
		defer db.Close()

		repo := repository.NewSqliteCartRepo(db)
		service := services.NewCartService(repo)

		// if repository.IsEmpty(items) {
		// 	color.Redln("userID:", userInfo.UserID, " Корзина пуста", err)
		// 	inlineKeyboard = messages.GetKeyboard("back", "Магазин")
		// 	messageText = "Корзина пуста"
		// 	botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		// 	session.LastBotMessageID = botMessageID
		// 	return
		// }

		messageText, err = service.GetCartText(int64(userInfo.UserID))
		if err != nil {
			color.Redln("userID:", userInfo.UserID, "Error:", err)
			inlineKeyboard = messages.GetKeyboard("back", "Магазин")
			messageText = "Произошла ошибка загрузки. Пожалуйста, попробуйте позже"
			botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
			session.LastBotMessageID = botMessageID
			return
		}

		inlineKeyboard = messages.GetKeyboard("Корзина", "Магазин")
		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	default:
		keyboard = messages.GetReplyKeyboard()
		messageText = messages.GetMessageText("start")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, keyboard)
	}

	SessionManager.PrintLogs(userInfo.UserID)
	session.LastBotMessageID = botMessageID
}

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	callbackMessageID := update.CallbackQuery.Message.MessageID
	userInfo := getUserInfo(update)

	session, exists := SessionManager.GetSession(userInfo.UserID)
	if !exists {
		session = SessionManager.CreateSession(userInfo)
	}
	session.LastUserMessageID = callbackMessageID

	if session.LastUserMessageID != 0 {
		messages.DeleteMessage(bot, session.LastUserMessageID, userInfo.UserID)
	}
	if session.LastBotMessageID != 0 {
		messages.DeleteMessage(bot, session.LastBotMessageID, userInfo.UserID)
	}

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var keyboard tgbotapi.ReplyKeyboardMarkup
	var messageText string
	var botMessageID int
	switch data {
	case "Магазин":
		inlineKeyboard = messages.GetKeyboard(data, nil)
		messageText = "Выберите категорию: "

		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "drinkware":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		service := initProductService(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := messages.GetKeyboard("back", "Магазин")
			messages.SendMessage(bot, userInfo.UserID, "Товаров нет.", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = card.GetTextTemplate()

		keyboard = messages.GetReplyKeyboard()
		messages.SendReplyKeyboard(bot, userInfo.UserID, "", keyboard)

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, card.Image)
	case "dishware":
		db, err := setupDatabase()
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		service := initProductService(db)

		products, err := service.GetProductByType(data)
		if err != nil {
			log.Println(err)
		}

		if repository.IsEmpty(products) {
			keyboard := messages.GetKeyboard("back", "Магазин")
			messages.SendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = card.GetTextTemplate()

		keyboard = messages.GetReplyKeyboard()
		messages.SendReplyKeyboard(bot, userInfo.UserID, "", keyboard)

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, card.Image)
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
			keyboard := messages.GetKeyboard("back", "Магазин")
			messages.SendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = card.GetTextTemplate()

		keyboard = messages.GetReplyKeyboard()
		messages.SendReplyKeyboard(bot, userInfo.UserID, "", keyboard)

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, card.Image)
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
			keyboard := messages.GetKeyboard("back", "Магазин")
			messages.SendMessage(bot, userInfo.UserID, "Товаров нет", keyboard)
			return
		}

		card := card.NewCard(products)
		session.CardManager.UpdateInfo(data, card)

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = card.GetTextTemplate()

		keyboard = messages.GetReplyKeyboard()
		messages.SendReplyKeyboard(bot, userInfo.UserID, "", keyboard)

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, card.Image)
	case "prev":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		color.Blueln("prev")
		currentCard.Prev()

		keyboard := messages.GetDynamicKeyboard("card", session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
	case "next":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		color.Blueln("next")
		currentCard.Next()

		keyboard := messages.GetDynamicKeyboard("card", session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
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

		total, err := service.AddItem(item)
		if err != nil {
			color.Redln("Error to add item:", err)
		}
		color.Redln("total in db:", total)

		session.CartManager.Add(item)
		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
	case "increment":
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
		}

		total, err := service.Increment(item)
		if err != nil {
			color.Redln("Error to increment item:", err)
		}
		session.CartManager.UpdateQuantity(currentCard.ID, total)
		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
	case "decrement":
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
		}

		total, err := service.Decrement(item)
		if err != nil {
			color.Redln("Error to increment item:", err)
		}
		color.Redln("Total in increment func:", total)

		session.CartManager.UpdateQuantity(currentCard.ID, total)
		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
	}

	SessionManager.PrintLogs(userInfo.UserID)
	session.LastBotMessageID = botMessageID
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