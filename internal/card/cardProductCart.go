package card

import (
	"fmt"

	"github.com/Resolution-hash/shop_bot/internal/repository"
)

type CardProductCart struct {
	ID                int
	Name              string
	Type              string
	Description       string
	Quantity          int
	Price             float64
	Image             string
	TotalCards        int
	CurrentCardNumber int
	ProductList       []*repository.CartProduct
}

func NewCardProductCart(cartProduct []*repository.CartProduct) *CardProductCart {
	if len(cartProduct) > 0 {
		return &CardProductCart{
			ID:                cartProduct[0].ProductID,
			Name:              cartProduct[0].Name,
			Description:       cartProduct[0].Description,
			Quantity:          cartProduct[0].Quantity,
			Price:             cartProduct[0].Price,
			Image:             cartProduct[0].Image,
			TotalCards:        len(cartProduct),
			CurrentCardNumber: 0,
			ProductList:       cartProduct,
		}
	}
	return nil
}

func (c *CardProductCart) Prev() {
	if c.CurrentCardNumber > 0 {
		c.CurrentCardNumber--
		c.ID = c.ProductList[c.CurrentCardNumber].ProductID
		c.Name = c.ProductList[c.CurrentCardNumber].Name
		c.Description = c.ProductList[c.CurrentCardNumber].Description
		c.Price = c.ProductList[c.CurrentCardNumber].Price
		c.Image = c.ProductList[c.CurrentCardNumber].Image
	}

}
func (c *CardProductCart) Next() {
	if c.CurrentCardNumber < c.TotalCards-1 {
		c.CurrentCardNumber++
		c.ID = c.ProductList[c.CurrentCardNumber].ProductID
		c.Name = c.ProductList[c.CurrentCardNumber].Name
		c.Description = c.ProductList[c.CurrentCardNumber].Description
		c.Price = c.ProductList[c.CurrentCardNumber].Price
		c.Image = c.ProductList[c.CurrentCardNumber].Image
	}

}

func (c *CardProductCart) GetTextTemplate() string {
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCardNumber+1, c.TotalCards)
}
