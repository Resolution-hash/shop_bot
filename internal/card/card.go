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
	Image       string
	TotalCards  int
	CurrentCardNumber int
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
	color.Yellowln("Image:", cm.CurrentCard.Image)
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
			Image:       products[0].Image,
			TotalCards:  len(products),
			CurrentCardNumber: 0,
			ProductList: products,
		}
	}
	return nil
}

func (c *Card) Prev() {
	if c.CurrentCardNumber > 0 {
		c.CurrentCardNumber--
		c.ID = c.ProductList[c.CurrentCardNumber].ID
		c.Type = c.ProductList[c.CurrentCardNumber].Type
		c.Name = c.ProductList[c.CurrentCardNumber].Name
		c.Description = c.ProductList[c.CurrentCardNumber].Description
		c.Price = c.ProductList[c.CurrentCardNumber].Price
		c.Image = c.ProductList[c.CurrentCardNumber].Image
	}

}
func (c *Card) Next() {
	if c.CurrentCardNumber < c.TotalCards-1 {
		c.CurrentCardNumber++
		c.ID = c.ProductList[c.CurrentCardNumber].ID
		c.Type = c.ProductList[c.CurrentCardNumber].Type
		c.Name = c.ProductList[c.CurrentCardNumber].Name
		c.Description = c.ProductList[c.CurrentCardNumber].Description
		c.Price = c.ProductList[c.CurrentCardNumber].Price
		c.Image = c.ProductList[c.CurrentCardNumber].Image
	}

}

func (c *Card) GetTextTemplate() string {
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCardNumber+1, c.TotalCards)
}
