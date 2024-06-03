package card

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/gookit/color"
)

type CartManager struct {
	Items map[int64]int
}

func NewCartManager() *CartManager {
	return &CartManager{
		Items: make(map[int64]int),
	}
}

func (c *CartManager) GetCartItemsDetails(userID int64) (string, error) {
	db, err := setupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := initCardService(db)

	details, err := service.GetCartInfo(int64(userID))
	if err != nil {
		return "", err
	}
	return details, nil

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
}

// func (c *CartManager) Get(item repository.CartItem) error {

// }

func (c *CartManager) Increment(item repository.CartItem) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		color.Redln("error to get cfg.DbUrl")
		return err
	}

	defer db.Close()

	service := initCardService(db)

	total, err := service.Increment(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

func (c *CartManager) Decrement(item repository.CartItem) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		color.Redln("error to get cfg.DbUrl")
		return err
	}

	defer db.Close()

	service := initCardService(db)

	total, err := service.Decrement(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

func (c *CartManager) AddToCart(item repository.CartItem) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		color.Redln("error to get cfg.DbUrl")
		return err
	}

	defer db.Close()

	service := initCardService(db)

	total, err := service.AddItem(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

// func (c *Cart) UpdateQuantity(productID int64, quantity int) {
// 	c.Items[productID] = quantity
// }

func (c *CartManager) Total(productID int64) string {
	return strconv.Itoa(c.Items[productID])
}

func (c *CartManager) PrintLogs() {
	fmt.Print("___________________\n\n")
	for k, v := range c.Items {
		color.Yellowf("ProductID:%v, quantity:%v\n", k, v)
	}
	fmt.Print("___________________\n\n")
}

func initCardService(db *sql.DB) *services.CartService {
	repo := repository.NewSqliteCartRepo(db)
	return services.NewCartService(repo)

}
