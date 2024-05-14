package card

import (
	"fmt"

	"github.com/Resolution-hash/shop_bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Card struct {
	Name        string
	Description string
	Price       float64
	TotalCards  int
	CurrentCard int
	ProductList []repository.Product
}

type CardSession struct {
	Card     *Card
	Keyboard tgbotapi.InlineKeyboardMarkup
}

func NewCard(products []repository.Product) *Card {
	if len(products) > 0 {
		return &Card{
			Name:        products[0].Name,
			Description: products[0].Description,
			Price:       products[0].Price,
			TotalCards:  len(products),
			CurrentCard: 0,
			ProductList: products,
		}
	}
	return nil
}

func (c *Card) Prev() {
	if c.CurrentCard > 0 {
		c.CurrentCard--
		c.Name = c.ProductList[c.CurrentCard].Name
		c.Description = c.ProductList[c.CurrentCard].Description
		c.Price = c.ProductList[c.CurrentCard].Price
	}
	fmt.Println("\n\n\n", c.CurrentCard+1, c.TotalCards, "\n\n\n")

}
func (c *Card) Next() {
	if c.CurrentCard < c.TotalCards-1 {
		c.CurrentCard++
		c.Name = c.ProductList[c.CurrentCard].Name
		c.Description = c.ProductList[c.CurrentCard].Description
		c.Price = c.ProductList[c.CurrentCard].Price
	}
	fmt.Println("\n\n\n", c.CurrentCard+1, c.TotalCards, "\n\n\n")

}

func (c *Card) GetTextTemplate() string {
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCard+1, c.TotalCards)
}
