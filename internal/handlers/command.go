package handlers

import (
	"errors"

	"github.com/Resolution-hash/shop_bot/internal/logic"
	"github.com/Resolution-hash/shop_bot/internal/messages"
	"github.com/Resolution-hash/shop_bot/internal/storage"
	product "github.com/Resolution-hash/shop_bot/repository/product"
	user "github.com/Resolution-hash/shop_bot/repository/user"

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
	messages.DeleteMessages(bot, *session, userInfo.UserID)

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
	default:
		handleAdminsAction(bot, update, session)
	}

	SessionManager.PrintLogs(userInfo.UserID)

}

func getUserInfo(update tgbotapi.Update) *user.User {
	if update.CallbackQuery != nil {
		usr := update.CallbackQuery.From
		return &user.User{
			UserID:     usr.ID,
			First_name: usr.FirstName,
			Last_name:  usr.LastName,
			User_name:  usr.UserName,
		}

	}

	usr := update.Message.From
	return &user.User{
		UserID:     usr.ID,
		First_name: usr.FirstName,
		Last_name:  usr.LastName,
		User_name:  usr.UserName,
	}
}

func SendError(bot *tgbotapi.BotAPI, s *sessions.Session, callback string, err error) {
	inlineKeyboard := messages.GetKeyboard("back", s, callback)
	messageText := err.Error()
	botMessageID := messages.SendMessage(bot, s.User.UserID, messageText, inlineKeyboard)
	s.User.SettingStep = ""
	s.LastBotMessageID = botMessageID
}

func handleAdminsAction(bot *tgbotapi.BotAPI, update tgbotapi.Update, session *sessions.Session) {
	switch session.User.SettingStep {
	case "uploadProduct":
		color.Blueln("uploadProduct")

		if update.Message.Photo != nil {
			photos := *update.Message.Photo
			photoSize := photos[len(photos)-1]

			product, err := logic.ParseProduct(update.Message.Caption)
			if err != nil {
				color.Redln(err)
				SendError(bot, session, "addItem", err)
				session.User.SettingStep = ""
				return
			}

			objName, err := UploadPhotos(bot, photoSize)
			if err != nil {
				color.Redln(err)
				SendError(bot, session, "addItem", err)
				session.User.SettingStep = ""
				return
			}
			product.Image = objName
			color.Redln("imageName:", product.Image)

			session.UpdateNewProduct(product)

			inlineKeyboard := messages.GetAdminKeyboard(session)
			messageText := logic.GetTestText(product)

			botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, session.User.UserID, messageText, inlineKeyboard, product.Image)
			if err != nil {
				color.Redln(err)
				SendError(bot, session, "addItem", err)
				session.User.SettingStep = ""
				return
			}
			session.UpdateSettingStep("confirmСhanges")
			session.LastBotMessageID = botMessageID
		} else {
			inlineKeyboard := messages.GetKeyboard("back", session, "adminPanel")
			messageText := "Фотография не добавлена"
			botMessageID := messages.SendMessage(bot, session.User.UserID, messageText, inlineKeyboard)
			session.LastBotMessageID = botMessageID
		}
	case "changePhoto":
		color.Blueln("changePhoto")
		if len(*update.Message.Photo) > 1 {
			photos := *update.Message.Photo
			photoSize := photos[len(photos)-1]

			objName, err := UploadPhotos(bot, photoSize)
			if err != nil {
				color.Redln(err)
				session.UpdateSettingStep("")
				SendError(bot, session, "adminPanel", err)
				return
			}
			color.Redln("changePhoto, photo is downloaded:", objName)
			color.Redln("imageName:", objName)

			currentCard := session.CardManager.CurrentCard
			product := product.Product{
				ID:    currentCard.ID,
				Image: objName,
			}

			err = session.CardManager.UpdateCardImage(product)
			if err != nil {
				session.UpdateSettingStep("")

				minioErr := storage.MinIORemovePhoto(objName)
				if minioErr != nil {
					errorMessage := err.Error() + minioErr.Error()
					color.Redln(errorMessage)
					SendError(bot, session, "adminPanel", errors.New(errorMessage))
					return
				}

				color.Redln(err)
				SendError(bot, session, "adminPanel", err)
				return
			}

			session.UpdateSettingStep("")
			inlineKeyboard := messages.GetKeyboard("back", session, "changeItem")
			messageText := "Изображения успешно изменено"
			botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, session.User.UserID, messageText, inlineKeyboard, objName)
			if err != nil {
				color.Redln(err)
				SendError(bot, session, "adminPanel", err)
				return
			}
			session.LastBotMessageID = botMessageID
		}
	case "сhangeText":
		session.UpdateSettingStep("")
		data, err := logic.ParseProduct(update.Message.Text)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "changeItem", err)
			return
		}
		currentCard := session.CardManager.CurrentCard
		product := product.Product{
			ID:          currentCard.ID,
			Name:        data.Name,
			Type:        data.Type,
			Description: data.Description,
			Price:       data.Price,
		}
		err = session.CardManager.UpdateCardText(product)
		if err != nil {
			color.Redln(err)
			SendError(bot, session, "changeItem", err)
			return
		}

		inlineKeyboard := messages.GetKeyboard("back", session, "changeItem")
		messageText := "Текст товара успешно изменен"
		botMessageID := messages.SendMessage(bot, session.User.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID

	default:
		inlineKeyboard := messages.GetKeyboard(update.Message.Text, session, nil)
		messageText := "Ошибка команды"
		botMessageID := messages.SendMessage(bot, session.User.UserID, messageText, inlineKeyboard)
		session.LastBotMessageID = botMessageID
	}

}
