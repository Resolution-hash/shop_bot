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
		} else if session.User.SettingStep != "" {
			keyboard = messages.GetAdminKeyboard(session)
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
		} else if session.User.SettingStep != "" {
			keyboard = messages.GetAdminKeyboard(session)
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
		session.UpdateSettingStep("adminPanel")
		messageText = "Панель администратора"
		inlineKeyboard := messages.GetAdminKeyboard(session)
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
		session.UpdateSettingStep("")
	case "addItem":
		session.UpdateSettingStep("uploadProduct")
		messageText := "Добавитьте фото и заполните карточку по примеру\nПример:\nКружка праздничная\ndrinkware\nОбъем: 450 мл\n599\n\nДоступные типы для карточке:\ndishware\ndrinkware\ncandles\n"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	case "confirmСhanges":
		color.Redln(session.NewProduct)
		err := logic.CreateProduct(session.NewProduct)
		if err != nil {
			session.UpdateNewProduct("")
			color.Redln(err)
			err := RemovePhotoFromStorage(session.NewProduct.Image)
			if err != nil {
				color.Redln(err)
			}
			SendError(bot, session, "addItem", err)
			session.LastBotMessageID = botMessageID
			return
		}
		session.UpdateNewProduct("")

		color.Greenln("product added")

		messageText := "Товар добавлен"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
		session.UpdateSettingStep("")

	case "cancelChanges":
		err := RemovePhotoFromStorage(session.NewProduct.Image)
		if err != nil {
			color.Redln(err)
			err := RemovePhotoFromStorage(session.NewProduct.Image)
			SendError(bot, session, "addItem", err)
			session.LastBotMessageID = botMessageID
			return
		}

		messageText := "Изменения отменены"
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
		session.UpdateSettingStep("")

	case "changeItem":
		session.UpdateSettingStep("changeItem")
		err := session.CardManager.GetCardAll(data)
		if err != nil {
			color.Redln(err)
		}

		inlineKeyboard := messages.GetAdminKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "adminPanel", err)
			session.LastBotMessageID = botMessageID
			session.UpdateSettingStep("")
			return
		}
		session.LastBotMessageID = botMessageID
	case "сhangePhoto":
		session.UpdateSettingStep("changePhoto")

		messageText := "Загрузите изображение формата jpeg."
		inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
		session.UpdateSettingStep("")
	case "сhangeText":
		session.UpdateSettingStep("сhangeText")

		messageText := "Заполните карточку по примеру\nПример:\nКружка праздничная\ndrinkware\nОбъем: 450 мл\n599\n\nДоступные типы для карточке:\ndishware\ndrinkware\ncandles\n"
		inlineKeyboard := messages.GetKeyboard("back", session, "changeItem")
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

		session.UpdateSettingStep("")
	case "deleteItems":
		session.UpdateSettingStep("deleteItems")
		err := session.CardManager.GetCardAll(data)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "adminPanel", err)
			session.LastBotMessageID = botMessageID
			session.UpdateSettingStep("")
			return
		}

		inlineKeyboard := messages.GetAdminKeyboard(session)
		messageText = session.CardManager.GetCardText()
		cardImage := session.CardManager.GetCardImage()

		botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, cardImage)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "adminPanel", err)
			session.LastBotMessageID = botMessageID
			session.UpdateSettingStep("")
			return
		}
		session.LastBotMessageID = botMessageID

	case "deleteProduct":
		session.UpdateSettingStep("")
		err := session.CardManager.DeleteCard(session.CardManager.CurrentCard.ID)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "deleteItems", err)
			session.LastBotMessageID = botMessageID
			return
		}
		inlineKeyboard := messages.GetKeyboard("back", session, "deleteItems")
		messageText := "Товар успешно удален"
		botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
		color.Yellowln(session.User.SettingStep)
	}
	SessionManager.PrintLogs(userInfo.UserID)
}
