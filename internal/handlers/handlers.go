package handlers

import (
	"github.com/Resolution-hash/shop_bot/internal/messages"
	"github.com/Resolution-hash/shop_bot/internal/repository"

	// "github.com/Resolution-hash/shop_bot/internal/services"
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
	deleteMessages(bot, *session, userInfo.UserID)

	var keyboard tgbotapi.ReplyKeyboardMarkup
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var botMessageID int
	var messageText string
	switch data {
	case "/start":
		keyboard = messages.GetReplyKeyboard()
		messageText = messages.GetMessageText("start")

		botMessageID = messages.SendReplyKeyboard(bot, userInfo.UserID, messageText, keyboard)

		session.LastBotMessageID = botMessageID

	case "Магазин":

		inlineKeyboard = messages.GetKeyboard(data, nil)
		messageText = "Выберите категорию: "
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	case "Корзина":

		// db, err := setupDatabase()
		// if err != nil {
		// 	color.Redln(err)
		// }
		// defer db.Close()

		// repo := repository.NewSqliteCartRepo(db)
		// service := services.NewCartService(repo)

		// if repository.IsEmpty(items) {
		// 	color.Redln("userID:", userInfo.UserID, " Корзина пуста", err)
		// 	inlineKeyboard = messages.GetKeyboard("back", "Магазин")
		// 	messageText = "Корзина пуста"
		// 	botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		// 	session.LastBotMessageID = botMessageID
		// 	return
		// }

		// messageText, err := service.GetCartInfo(int64(userInfo.UserID))
		// if err != nil {
		// 	color.Redln("userID:", userInfo.UserID, "Error:", err)
		// 	inlineKeyboard = messages.GetKeyboard("back", "Магазин")
		// 	messageText = "Произошла ошибка загрузки. Пожалуйста, попробуйте позже"
		// 	botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		// 	session.LastBotMessageID = botMessageID
		// 	return
		// }

		inlineKeyboard = messages.GetKeyboard("buttonForCart", "Магазин")
		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	default:
		keyboard = messages.GetReplyKeyboard()
		messageText = messages.GetMessageText("start")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, keyboard)
		session.LastBotMessageID = botMessageID
	}

	SessionManager.PrintLogs(userInfo.UserID)

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
	deleteMessages(bot, *session, userInfo.UserID)

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var messageText string
	var botMessageID int
	switch data {
	case "Магазин":
		inlineKeyboard = messages.GetKeyboard(data, nil)
		messageText = "Выберите категорию: "

		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "drinkware":
		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "dishware":

		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "drawings":
		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "showAllItems":
		err := session.CardManager.GetCardAll(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	// case "changeCart":
	// 	db, err := setupDatabase()
	// 	if err != nil {
	// 		color.Redln(err)
	// 	}
	// 	defer db.Close()

	// 	repo := repository.NewSqliteCartRepo(db)
	// 	service := services.NewCartService(repo)

	// 	cartProducts, err := service.GetItemsByUserID(int64(userInfo.UserID))
	// 	if err != nil {
	// 		keyboard := messages.GetKeyboard("back", "Магазин")
	// 		messages.SendMessage(bot, userInfo.UserID, "Ошибка загрузки корзины. Пожалуйста, попробуйте позже", keyboard)
	// 		return
	// 	}

	// 	if repository.IsEmpty(cartProducts) {
	// 		keyboard := messages.GetKeyboard("back", "Магазин")
	// 		messages.SendMessage(bot, userInfo.UserID, "Корзина пуста", keyboard)
	// 		return
	// 	}

	// 	card := card.NewCardProductCart(cartProducts)
	// 	session.CardManager.UpdateInfoCart(data, card)

	// 	inlineKeyboard := messages.GetCartKeyboard(session)
	// 	messageText = card.GetTextTemplate()

	// 	botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, card.Image)
	// 	session.LastBotMessageID = botMessageID

	case "prev":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		color.Blueln("prev")
		session.CardManager.PrevCard()

		keyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "next":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		color.Blueln("next")
		session.CardManager.NextCard()

		keyboard := messages.GetDynamicKeyboard("card", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "prevProductCart":
		// defer func(s *sessions.Session) {
		// 	s.CardManager.PrintLogs()
		// 	s.CartManager.PrintLogs()
		// }(session)

		currentCard := session.CardManager.CurrentProductCart

		color.Blueln("prevProductCart")
		currentCard.Prev()

		keyboard := messages.GetCartKeyboard(session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
		session.LastBotMessageID = botMessageID
	case "nextProductCart":
		// defer func(s *sessions.Session) {
		// 	s.CardManager.PrintLogs()
		// 	s.CartManager.PrintLogs()
		// }(session)

		currentCard := session.CardManager.CurrentProductCart

		color.Blueln("nextProductCart")
		currentCard.Next()

		keyboard := messages.GetCartKeyboard(session)
		messageText = currentCard.GetTextTemplate()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
		session.LastBotMessageID = botMessageID

	case "addToCart":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			// s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		item := repository.CartItem{
			ProductID: currentCard.ID,
			UserID:    int64(userInfo.UserID),
			Quantity:  1,
		}

		err := session.CartManager.AddToCart(item)
		if err != nil {
			color.Redln(err)
		}

		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "increment":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		item := repository.CartItem{
			ProductID: currentCard.ID,
			UserID:    int64(userInfo.UserID),
		}

		err := session.CartManager.Increment(item)
		if err != nil {
			color.Redln(err)
		}

		// db, err := setupDatabase()
		// if err != nil {
		// 	color.Redln(err)
		// }
		// defer db.Close()

		// repo := repository.NewSqliteCartRepo(db)
		// service := services.NewCartService(repo)

		// item := repository.CartItem{
		// 	ProductID: currentCard.ID,
		// 	UserID:    int64(userInfo.UserID),
		// }

		// total, err := service.Increment(item)
		// if err != nil {
		// 	color.Redln("Error to increment item:", err)
		// }
		// session.CartManager.UpdateQuantity(currentCard.ID, total)
		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "decrement":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		item := repository.CartItem{
			ProductID: currentCard.ID,
			UserID:    int64(userInfo.UserID),
		}

		err := session.CartManager.Decrement(item)
		if err != nil {
			color.Redln(err)
		}

		// db, err := setupDatabase()
		// if err != nil {
		// 	color.Redln(err)
		// }
		// defer db.Close()

		// repo := repository.NewSqliteCartRepo(db)
		// service := services.NewCartService(repo)

		// item := repository.CartItem{
		// 	ProductID: currentCard.ID,
		// 	UserID:    int64(userInfo.UserID),
		// }

		// total, err := service.Decrement(item)
		// if err != nil {
		// 	color.Redln("Error to increment item:", err)
		// }
		// color.Redln("Total increment func:", total)

		// session.CartManager.UpdateQuantity(currentCard.ID, total)
		keyboard := messages.GetDynamicKeyboard("addToCart", session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "incrementProductCart":
		// defer func(s *sessions.Session) {
		// 	s.CardManager.PrintLogs()
		// 	s.CartManager.PrintLogs()
		// }(session)

		// currentCard := session.CardManager.CurrentProductCart

		// db, err := setupDatabase()
		// if err != nil {
		// 	color.Redln(err)
		// }
		// defer db.Close()

		// repo := repository.NewSqliteCartRepo(db)
		// service := services.NewCartService(repo)

		// item := repository.CartItem{
		// 	ProductID: int64(currentCard.ID),
		// 	UserID:    int64(userInfo.UserID),
		// }

		// total, err := service.Increment(item)
		// if err != nil {
		// 	color.Redln("Error to increment item:", err)
		// }
		// session.CartManager.UpdateQuantity(int64(currentCard.ID), total)
		// keyboard := messages.GetCartKeyboard(session)
		// messageText = currentCard.GetTextTemplate()

		// botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
		session.LastBotMessageID = botMessageID
		// case "decrementProductCart":
		// 	// defer func(s *sessions.Session) {
		// 	// 	s.CardManager.PrintLogs()
		// 	// 	s.CartManager.PrintLogs()
		// 	// }(session)

		// 	currentCard := session.CardManager.CurrentProductCart

		// 	db, err := setupDatabase()
		// 	if err != nil {
		// 		color.Redln(err)
		// 	}
		// 	defer db.Close()

		// 	repo := repository.NewSqliteCartRepo(db)
		// 	service := services.NewCartService(repo)

		// 	item := repository.CartItem{
		// 		ProductID: int64(currentCard.ID),
		// 		UserID:    int64(userInfo.UserID),
		// 	}

		// 	total, err := service.Decrement(item)
		// 	if err != nil {
		// 		color.Redln("Error to increment item:", err)
		// 	}
		// 	color.Redln("Total increment func:", total)

		// 	session.CartManager.UpdateQuantity(int64(currentCard.ID), total)
		// 	keyboard := messages.GetCartKeyboard(session)
		// 	messageText = currentCard.GetTextTemplate()

		// 	botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, currentCard.Image)
		// 	session.LastBotMessageID = botMessageID
	}

	SessionManager.PrintLogs(userInfo.UserID)
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

func deleteMessages(bot *tgbotapi.BotAPI, session sessions.Session, userId int) {
	color.Redln("LastUserMessageID", session.LastUserMessageID)
	color.Redln("LastBotMessageID", session.LastBotMessageID)

	if session.LastUserMessageID != 0 && session.LastUserMessageID != session.LastBotMessageID {
		color.Redln("delete user message", session.LastUserMessageID)

		messages.DeleteMessage(bot, session.LastUserMessageID, userId)
		session.LastUserMessageID = 0
	}
	if session.LastBotMessageID != 0 && session.LastUserMessageID != session.LastBotMessageID {
		color.Redln("delete bot message", session.LastBotMessageID)

		messages.DeleteMessage(bot, session.LastBotMessageID, userId)
		session.LastBotMessageID = 0
	}
}
