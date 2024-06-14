package handlers

import (
	"github.com/Resolution-hash/shop_bot/internal/logic"
	"github.com/Resolution-hash/shop_bot/internal/messages"
	"github.com/Resolution-hash/shop_bot/internal/sessions"
	cart "github.com/Resolution-hash/shop_bot/repository/cart"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	callbackMessageID := update.CallbackQuery.Message.MessageID
	userInfo := getUserInfo(update)

	session, exists := SessionManager.GetSession(userInfo.UserID)
	if !exists {
		session = SessionManager.CreateSession(userInfo)
	}
	session.LastUserMessageID = callbackMessageID
	messages.DeleteMessages(bot, *session, userInfo.UserID)

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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
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

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.LastBotMessageID = botMessageID
	case "addToCart":
		session.UpdateStep(data)
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard

		item := cart.CartItem{
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

		botMessageID, err = messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.LastBotMessageID = botMessageID
	case "increment":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard
		item := cart.CartItem{
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

		botMessageID, err = messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.LastBotMessageID = botMessageID
	case "decrement":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard
		item := cart.CartItem{
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

		botMessageID, err = messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.LastBotMessageID = botMessageID
	case "delete":
		defer func(s *sessions.Session) {
			s.CardManager.PrintLogs()
			s.CartManager.PrintLogs()
		}(session)

		currentCard := session.CardManager.CurrentCard
		item := cart.CartItem{
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

		botMessageID, err = messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, keyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, session.PrevStep, err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.LastBotMessageID = botMessageID

	//ADMINS____
	case "adminPanel":
		messageText = "Панель администратора"
		inlineKeyboard := messages.GetAdminKeyboard(session)
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "addItem":
		session.UpdateSettingStep("uploadProduct")
		messageText := "Последоватальность действий для добавления карточки:\n1. Добавить изображение в MinIo(файловое хранилище)\n2. Написать данные для карточки (Название,тип,описание,цена,название загруженной картинка в MinIO)\n\nПример:\n\nКружка праздничная\ndrinkware\nОбъем: 450 мл\n599\ncup-4\n\nВАЖНО\n\nДоступные типы:\ndishware\ndrinkware\ncandles\n\nЕсли указать неправильный тип, карточка не появится в разделе"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "confirmСhanges":
		color.Redln(session.TestProduct)
		err := logic.CreateProduct(session.TestProduct)
		if err != nil {
			session.UpdateTestProduct("")
			color.Redln(err)
			err := RemovePhotoFromStorage(session.TestProduct.Image)
			if err != nil {
				color.Redln(err)
			}
			SendError(bot, session, "addItem", err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.UpdateTestProduct("")

		color.Greenln("product added")

		messageText := "Товар добавлен"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.UpdateSettingStep("")
		session.LastBotMessageID = botMessageID
	case "cancelChanges":
		err := RemovePhotoFromStorage(session.TestProduct.Image)
		if err != nil {
			color.Redln(err)
			err := RemovePhotoFromStorage(session.TestProduct.Image)
			SendError(bot, session, "addItem", err)
			session.LastBotMessageID = botMessageID
			return
		}

		messageText := "Товар добавлен"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.UpdateSettingStep("")
		session.LastBotMessageID = botMessageID
	}
	SessionManager.PrintLogs(userInfo.UserID)
}
