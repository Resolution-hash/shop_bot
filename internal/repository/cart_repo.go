package repository

import (
	"fmt"
	"strconv"

	"github.com/gookit/color"
)

type CartItem struct {
	ProductID int64
	UserID    int64
	Quantity  int
}

type CartProduct struct {
	ProductID   int
	Name        string
	Description string
	Price       float64
	Quantity    int
}

type Cart struct {
	Items map[int64]int
}

func NewCart() *Cart {
	return &Cart{
		Items: make(map[int64]int),
	}
}

func (c *Cart) Add(item CartItem) {
	c.Items[item.ProductID] = item.Quantity
}

func (c *Cart) UpdateQuantity(productID int64, quantity int) {
	c.Items[productID] = quantity
}

func (c *Cart) Total(productID int64) string {
	return strconv.Itoa(c.Items[productID])
}

func (c *Cart) PrintLogs() {
	fmt.Print("___________________\n\n")
	for k, v := range c.Items {
		color.Yellowf("ProductID:%v, quantity:%v\n", k, v)
	}
	fmt.Print("___________________\n\n")
}

type CartRepo interface {
	AddItem(CartItem) (int, error)
	Increment(CartItem) (int, error)
	Decrement(CartItem) (int, error)
	// RemoveItem(int64) error
	GetItemsByUserID(int64) ([]*CartProduct, error)
}