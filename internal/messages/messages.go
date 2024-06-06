package messages

import (
	"fmt"
	"log"
	"os"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/sessions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

func SendMessage(bot *tgbotapi.BotAPI, userID int, text string, keyboard interface{}) int {
	msg := tgbotapi.NewMessage(int64(userID), text)
	if keyboard != nil {
		switch k := keyboard.(type) {
		case tgbotapi.InlineKeyboardMarkup:
			msg.ReplyMarkup = k
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = k
		}
	}
	msg.ParseMode = "HTML"
	sentMsg, err := bot.Send(msg)
	if err != nil {
		color.Redln("Ошибка отправки сообщения: %s\n", err)
		return 0
	}

	return sentMsg.MessageID
}

func SendReplyKeyboard(bot *tgbotapi.BotAPI, userID int, text string, keyboard tgbotapi.ReplyKeyboardMarkup) int {
	msg := tgbotapi.NewMessage(int64(userID), text)
	msg.ReplyMarkup = keyboard
	sentMsg, err := bot.Send(msg)
	if err != nil {
		color.Redln("Ошибка отправки сообщения: %s", err)
		return 0
	}
	color.Greenln("Keyboard is fetched")
	return sentMsg.MessageID
}

func EditMessage(bot *tgbotapi.BotAPI, userID int, messageID int, text string) int {
	msg := tgbotapi.NewEditMessageText(int64(userID), messageID, "Новый текст")
	sentMsg, err := bot.Send(msg)
	if err != nil {
		color.Redln("Ошибка отправки сообщения: %s\n", err)
		return 0
	}
	color.Greenln("Keyboard is fetched")
	return sentMsg.MessageID
}

func SendMessageWithPhoto(bot *tgbotapi.BotAPI, userID int, text string, keyboard interface{}, imageName string) int {
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

	msg := tgbotapi.NewPhotoUpload(int64(userID), file.Name())
	msg.File = tgbotapi.FileReader{
		Name:   file.Name(),
		Reader: file,
		Size:   -1,
	}
	msg.Caption = text

	if keyboard != nil {
		switch k := keyboard.(type) {
		case tgbotapi.InlineKeyboardMarkup:
			msg.ReplyMarkup = k
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = k
		}
	}
	sentMsg, err := bot.Send(msg)
	if err != nil {
		color.Redln("Ошибка отправки сообщения:", err)
		return 0
	}
	return sentMsg.MessageID
}

// func SendMessageWithPhotos(bot *tgbotapi.BotAPI, userID int, text string, keyboard interface{}, imageNames []string) int {
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	color.Redln("imageNames LEN", len(imageNames))
// 	files := make(map[string]tgbotapi.RequestFileData)
// 	mediaGroup := make([]interface{}, len(imageNames))
// 	for i, imageName := range imageNames {
// 		color.Redln("imageName", imageName)
// 		path := cfg.ImagesUrl + "\\" + imageName + ".jpg"
// 		color.Redln(path)

// 		file, err := os.Open(path)
// 		if err != nil {
// 			fmt.Println("Error to upload file")
// 		}
// 		defer file.Close()

// 		photo := tgbotapi.NewInputMediaPhoto(path)

// 		if i == 0 {
// 			photo.Caption = text
// 		}

// 		mediaGroup[i] = photo
// 	}
// 	mediaGroupConfig := tgbotapi.NewMediaGroup(int64(userID), mediaGroup)

// 	if keyboard != nil {
// 		switch k := keyboard.(type) {
// 		case tgbotapi.InlineKeyboardMarkup:
// 			mediaGroupConfig.ReplyMarkup = k
// 		case tgbotapi.ReplyKeyboardMarkup:
// 			mediaGroupConfig.ReplyMarkup = k
// 		}
// 	}
// 	sentMsg, err := bot.Send(mediaGroupConfig)
// 	if err != nil {
// 		color.Redln("Ошибка отправки сообщения:", err)
// 		return 0
// 	}
// 	return sentMsg.MessageID
// }

func DeleteMessage(bot *tgbotapi.BotAPI, messageID int, userID int) {
	deleteConfig := tgbotapi.NewDeleteMessage(int64(userID), messageID)
	if _, err := bot.Send(deleteConfig); err != nil {
		log.Printf("Ошибка при удалении сообщения: %s\n", err)
	}
}

func GetReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Магазин"),
			tgbotapi.NewKeyboardButton("Корзина"),
		),
	)
	return keyboard
}

func GetKeyboard(value string, session *sessions.Session, back interface{}) tgbotapi.InlineKeyboardMarkup {
	switch value {

	case "Магазин":
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Свечи 🕯️", "candles"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Посуда для питья 🍷", "drinkware"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Посуда для еды 🍽️", "dishware"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать все 🔍", "showAllItems"),
			),
		)
	case "Корзина":
		if session.CartManager.CartIsEmpty {
			return tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
				),
			)
		}
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Изменить корзину ✏️", "changeCart"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Оформить заказ 📦", "Checkout"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
			),
		)
	case "back":
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", back.(string)),
			),
		)
	default:
		color.Redln("Value is not found in GetKeyboard()")
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
			),
		)
	}
}

func GetCardKeyboard(session *sessions.Session) tgbotapi.InlineKeyboardMarkup {
	itemID := int(session.CardManager.CurrentCard.ID)
	userID := session.User.UserID

	quantity, err := session.CartManager.GetQuantity(itemID, userID)
	if err != nil {
		color.Redln(err)
	}
	color.Redln("quantity", quantity, " itemID", itemID)

	if quantity != "0" {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏪", "prev"),
				tgbotapi.NewInlineKeyboardButtonData("⏩", "next"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➖", "decrement"),
				tgbotapi.NewInlineKeyboardButtonData(quantity, "no_action"),
				tgbotapi.NewInlineKeyboardButtonData("➕", "increment"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Перейти в корзину 🛒", "Корзина"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
			),
		)

	} else {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏪", "prev"),
				tgbotapi.NewInlineKeyboardButtonData("⏩", "next"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину  🛒", "addToCart"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
			),
		)
	}
}

// func GetDynamicKeyboard(value string, session *sessions.Session) tgbotapi.InlineKeyboardMarkup {
// 	itemID := int(session.CardManager.CurrentCard.ID)
// 	userID := session.User.UserID

// 	quantity, err := session.CartManager.GetQuantity(itemID, userID)

// 	color.Redln("total in getDynamicKeyboard func:", quantity)
// 	var cartButtons []tgbotapi.InlineKeyboardButton

// 	if total != "0" {
// 		cartButtons = tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("➖", "decrement"),
// 			tgbotapi.NewInlineKeyboardButtonData(quantity, "no_action"),
// 			tgbotapi.NewInlineKeyboardButtonData("➕", "increment"),
// 		)
// 	} else {
// 		cartButtons = tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину  🛒", "addToCart"),
// 		)

// 	}

// 	switch value {
// 	case "addToCart":
// 		keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 			tgbotapi.NewInlineKeyboardRow(
// 				tgbotapi.NewInlineKeyboardButtonData("⏪", "prev"),
// 				tgbotapi.NewInlineKeyboardButtonData("⏩", "next"),
// 			),
// 			cartButtons,
// 			tgbotapi.NewInlineKeyboardRow(
// 				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
// 			),
// 		)
// 		return keyboard
// 	case "card":
// 		keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 			tgbotapi.NewInlineKeyboardRow(
// 				tgbotapi.NewInlineKeyboardButtonData("⏪", "prev"),
// 				tgbotapi.NewInlineKeyboardButtonData("⏩", "next"),
// 			),
// 			cartButtons,
// 			tgbotapi.NewInlineKeyboardRow(
// 				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
// 			),
// 		)
// 		return keyboard

// 	case "back":
// 		keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 			tgbotapi.NewInlineKeyboardRow(
// 				tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
// 			),
// 		)
// 		return keyboard
// 	default:
// 		keyboard := tgbotapi.NewInlineKeyboardMarkup()
// 		color.Redln("value is not found on func getKeyboard()")
// 		return keyboard
// 	}

// }

// func GetCartKeyboard(session *sessions.Session) tgbotapi.InlineKeyboardMarkup {
// 	currentCard := session.CardManager.CurrentProductCart
// 	total := session.CartManager.Total(int64(currentCard.ID))
// 	color.Redln("total in getDynamicKeyboard func:", total)

// 	keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("⏪", "prevProductCart"),
// 			tgbotapi.NewInlineKeyboardButtonData("⏩", "nextProductCart"),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("➖", "decrementProductCart"),
// 			tgbotapi.NewInlineKeyboardButtonData(total, "no_action"),
// 			tgbotapi.NewInlineKeyboardButtonData("➕", "incrementProductCart"),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("Удалить из корзины", "Корзина"),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("Вернуться ⬅️", "Магазин"),
// 		),
// 	)
// 	return keyboard
// }

func GetMessageText(step string) string {
	switch step {
	case "start":
		return "🎉 Добро пожаловать в наш Магазин Керамики! \n\n🎨Мы рады видеть вас среди наших ценителей уникальной керамической продукции. Здесь вы найдете изысканные изделия, созданные для того, чтобы добавить уюта и красоты вашему дому."
	default:
		return "Такой команды нет. Пожалуйста, выберите из доступных команд"
	}
}
