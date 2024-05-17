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

func (c *Cart) Quantitiy(item CartItem) string {
	return strconv.Itoa(c.Items[item.ProductID])
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
	// Increment(int64) error
	// Decrement(int64) error
	// RemoveItem(int64) error
	// GetItemsByID(int64) ([]CartItem, error)
}
