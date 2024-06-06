package handlers

import (
	"github.com/Resolution-hash/shop_bot/internal/messages"
	"github.com/Resolution-hash/shop_bot/internal/repository"

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
	if data != "/start" {
		session.LastUserMessageID = userMessageID
	}
	deleteMessages(bot, *session, userInfo.UserID)

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var botMessageID int
	var messageText string
	switch data {
	case "/start":
		inlineKeyboard = messages.GetKeyboard("start", session, nil)
		messageText = messages.GetMessageText("start")

		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)

		session.LastBotMessageID = botMessageID

	case "Магазин":
		inlineKeyboard = messages.GetKeyboard(data, session, nil)
		messageText = "Выберите категорию: "
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	case "Корзина":
		messageText, err := session.CartManager.GetCartItemsDetails(int64(userInfo.UserID))
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard = messages.GetKeyboard(data, session, nil)
		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
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
		session.UpdateStep(data)
		inlineKeyboard = messages.GetKeyboard(data, session, nil)
		messageText = "Выберите категорию: "

		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "Корзина":
		session.UpdateStep(data)
		messageText, err := session.CartManager.GetCartItemsDetails(int64(userInfo.UserID))
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard = messages.GetKeyboard(data, session, nil)
		botMessageID := messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	case "drinkware":
		session.UpdateStep(data)
		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetCardKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "dishware":
		session.UpdateStep(data)
		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetCardKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "drawings":
		session.UpdateStep(data)
		err := session.CardManager.GetCardByType(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetCardKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "showAllItems":
		session.UpdateStep(data)
		err := session.CardManager.GetCardAll(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetCardKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "changeCart":
		session.UpdateStep(data)
		err := session.CardManager.GetCartItemsByUserID(data, userInfo.UserID)
		if err != nil {
			color.Redln(err)
		}
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()
		inlineKeyboard := messages.GetKeyboard(data, session, nil)

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "prev":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		color.Blueln("prev")
		session.CardManager.PrevCard()

		prevStep := session.PrevStep
		var keyboard tgbotapi.InlineKeyboardMarkup
		if prevStep == "Корзина" {
			keyboard = messages.GetKeyboard("changeCart", session, nil)
		} else {
			keyboard = messages.GetCardKeyboard(session)
		}
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

		prevStep := session.PrevStep
		var keyboard tgbotapi.InlineKeyboardMarkup
		if prevStep == "Корзина" {
			keyboard = messages.GetKeyboard("changeCart", session, nil)
		} else {
			keyboard = messages.GetCardKeyboard(session)
		}
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "addToCart":
		session.UpdateStep(data)
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

		keyboard := messages.GetCardKeyboard(session)
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

		prevStep := session.PrevStep
		var keyboard tgbotapi.InlineKeyboardMarkup
		if prevStep == "Корзина" {
			keyboard = messages.GetKeyboard("changeCart", session, nil)
		} else {
			keyboard = messages.GetCardKeyboard(session)
		}

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

		prevStep := session.PrevStep
		var keyboard tgbotapi.InlineKeyboardMarkup

		if prevStep == "Корзина" {
			keyboard = messages.GetKeyboard("changeCart", session, nil)
		} else {
			keyboard = messages.GetCardKeyboard(session)
		}

		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID
	case "delete":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard
		item := repository.CartItem{
			ProductID: currentCard.ID,
			UserID:    int64(userInfo.UserID),
			Quantity:  0,
		}

		err := session.CartManager.DeleteItem(item)
		if err != nil {
			color.Redln("Error delete item", err)
		}

		prevStep := session.PrevStep
		var keyboard tgbotapi.InlineKeyboardMarkup
		if prevStep == "Корзина" {
			keyboard = messages.GetKeyboard("changeCart", session, nil)
		} else {
			keyboard = messages.GetCardKeyboard(session)
		}

		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID = messages.SendMessageWithPhoto(bot, userInfo.UserID, messageText, keyboard, cardImage)
		session.LastBotMessageID = botMessageID

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

	if session.LastUserMessageID != 0 {
		color.Redln("delete user message", session.LastUserMessageID)

		messages.DeleteMessage(bot, session.LastUserMessageID, userId)
		session.LastUserMessageID = 0
	}
	if session.LastBotMessageID != 0 {
		color.Redln("delete bot message", session.LastBotMessageID)

		messages.DeleteMessage(bot, session.LastBotMessageID, userId)
		session.LastBotMessageID = 0
	}
}
