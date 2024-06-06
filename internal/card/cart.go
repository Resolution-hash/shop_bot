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
	Items       map[int64]int
	CartIsEmpty bool
}

func NewCartManager() *CartManager {
	return &CartManager{
		Items:       make(map[int64]int),
		CartIsEmpty: true,
	}
}

func (c *CartManager) ChangeCartStatus(status bool) {
	c.CartIsEmpty = status
}

func (c *CartManager) GetCartItemsDetails(userID int64) (string, error) {
	db, err := setupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCardService(db)

	products, err := service.GetItemsByUserID(userID)
	if err != nil {
		return "", err
	}
	if len(products) == 0 {
		c.ChangeCartStatus(true)
		return "В корзине еще нет товаров.", nil
	}
	c.ChangeCartStatus(false)

	return service.FormatCartText(products), nil
}

// func (c *CartManager) GetItemsByUserID(userID int64) ([]*repository.CartProduct, error) {
// 	db, err := setupDatabase()
// 	if err != nil {
// 		color.Redln(err)
// 	}
// 	defer db.Close()

// 	service := InitCardService(db)

// 	products, err := service.GetItemsByUserID(int64(userID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return products, nil

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

	service := InitCardService(db)

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

	service := InitCardService(db)

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

	service := InitCardService(db)

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

// func (c *CartManager) Total(productID int64) string {
// 	return strconv.Itoa(c.Items[productID])
// }

func (c *CartManager) GetQuantity(itemID int, userID int) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		color.Redln("error to get cfg.DbUrl")
		return "", err
	}

	defer db.Close()

	service := InitCardService(db)

	item := repository.CartItem{
		ProductID: int64(itemID),
		UserID:    int64(userID),
		Quantity:  0,
	}

	total, err := service.GetQuantityByItemID(item)
	if err != nil {
		color.Redln("Error to GetQuantityByItemID", err)
	}
	color.Redln("total", total)

	return strconv.Itoa(total), nil
}

func (c *CartManager) PrintLogs() {
	fmt.Print("___________________\n\n")
	for k, v := range c.Items {
		color.Yellowf("ProductID:%v, quantity:%v\n", k, v)
	}
	fmt.Print("___________________\n\n")
}

func InitCardService(db *sql.DB) *services.CartService {
	repo := repository.NewSqliteCartRepo(db)
	return services.NewCartService(repo)
}
