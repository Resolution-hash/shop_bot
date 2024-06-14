package handlers

import (
	"github.com/Resolution-hash/shop_bot/internal/logic"
	"github.com/Resolution-hash/shop_bot/internal/messages"
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
		if session.User.SettingStep == "uploadProduct" {
			color.Blueln("uploadProduct")

			//Проверить, есть ли фото. Если есть, загрузить.
			//Проверить данные карточки.
			//Отправить сообщение с этими данными для подтверждения

			if len(*update.Message.Photo) > 1 {
				photos := *update.Message.Photo
				photoSize := photos[len(photos)-1]

				product, err := logic.ParseProduct(update.Message.Caption)
				if err != nil {
					color.Redln(err)
					SendError(bot, session, "addItem", err)
					return
				}

				objName, err := UploadPhotos(bot, photoSize)
				if err != nil {
					color.Redln(err)
					SendError(bot, session, "addItem", err)
					return
				}
				product.Image = objName
				color.Redln("imageName:", product.Image)

				session.UpdateTestProduct(product)

				inlineKeyboard := messages.GetAdminCardSetting()
				messageText = logic.GetTestText(product)

				botMessageID, err := messages.SendMessageWithPhotoMinIO(bot, userInfo.UserID, messageText, inlineKeyboard, product.Image)
				if err != nil {
					color.Redln(err)
					SendError(bot, session, "addItem", err)
					session.User.SettingStep = ""
					return
				}
				session.UpdateSettingStep("confirmСhanges")
				session.LastBotMessageID = botMessageID
			}

		} else {
			inlineKeyboard = messages.GetKeyboard(data, session, nil)
			messageText = "Ошибка команды"
			botMessageID = messages.SendMessage(bot, userInfo.UserID, messageText, inlineKeyboard)
			session.LastBotMessageID = botMessageID
		}

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
