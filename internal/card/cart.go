package card

import (
	"database/sql"
	"fmt"
	"strconv"

	cart "github.com/Resolution-hash/shop_bot/internal/repository/cart"
	db "github.com/Resolution-hash/shop_bot/internal/repository"

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

func (c *CartManager) DeleteItem(item cart.CartItem) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	err = service.DeleteItem(item)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartManager) GetCartItemsDetails(userID int64) (string, error) {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

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

func (c *CartManager) Increment(item cart.CartItem) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	total, err := service.Increment(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

func (c *CartManager) Decrement(item cart.CartItem) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	total, err := service.Decrement(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

func (c *CartManager) AddToCart(item cart.CartItem) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	total, err := service.AddItem(item)
	if err != nil {
		color.Redln("Error to increment item:", err)
	}

	c.Items[item.ProductID] = total
	return nil
}

func (c *CartManager) GetQuantity(itemID int, userID int) (string, error) {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	item := cart.CartItem{
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

func InitCartService(db *sql.DB) *services.CartService {
	repo := cart.NewSqliteCartRepo(db)
	return services.NewCartService(repo)
}


