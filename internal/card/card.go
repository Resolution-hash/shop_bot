package card

import (
	"fmt"

	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/gookit/color"
)

type Card struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Price       float64
	TotalCards  int
	CurrentCard int
	ProductList []repository.Product
}

type CardManager struct {
	Cards       map[string]*Card
	CurrentCard *Card
}

func NewCardManager() *CardManager {
	return &CardManager{
		Cards: map[string]*Card{},
	}
}

// func (cm *CardManager) GetCard(data string) *Card {
// 	card, exists := cm.Cards[data]
// 	if exists {
// 		return card
// 	}
// 	return nil
// }

func (cm *CardManager) UpdateInfo(step string, card *Card) {
	cm.Cards[step] = card
	cm.CurrentCard = card
}

func (cm *CardManager) PrintLogs() {
	fmt.Print("___________________\n\n")
	color.Yellowln("ProductID:", cm.CurrentCard.ID)
	color.Yellowln("ProductName:", cm.CurrentCard.Name)
	color.Yellowln("ProductType:", cm.CurrentCard.Type)
	color.Yellowln("ProductDescription:", cm.CurrentCard.Description)
	color.Yellowln("ProductPrice:", cm.CurrentCard.Price)
	fmt.Print("___________________\n\n")
}

func NewCard(products []repository.Product) *Card {
	if len(products) > 0 {
		return &Card{
			ID:          products[0].ID,
			Name:        products[0].Name,
			Type:        products[0].Type,
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
		c.ID = c.ProductList[c.CurrentCard].ID
		c.Type = c.ProductList[c.CurrentCard].Type
		c.Name = c.ProductList[c.CurrentCard].Name
		c.Description = c.ProductList[c.CurrentCard].Description
		c.Price = c.ProductList[c.CurrentCard].Price
	}

}
func (c *Card) Next() {
	if c.CurrentCard < c.TotalCards-1 {
		c.CurrentCard++
		c.ID = c.ProductList[c.CurrentCard].ID
		c.Type = c.ProductList[c.CurrentCard].Type
		c.Name = c.ProductList[c.CurrentCard].Name
		c.Description = c.ProductList[c.CurrentCard].Description
		c.Price = c.ProductList[c.CurrentCard].Price
	}

}

func (c *Card) GetTextTemplate() string {
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCard+1, c.TotalCards)
}
