package logic

import (
	"fmt"
	"log"

	db "github.com/Resolution-hash/shop_bot/repository/db"

	"github.com/gookit/color"
)

type ProductInfo interface {
	GetID() int64
	GetName() string
	GetType() string
	GetDescription() string
	GetPrice() float64
	GetImage() string
}

type Card struct {
	ID                int64
	Name              string
	Type              string
	Description       string
	Quantity          int
	Price             float64
	Image             string
	TotalCards        int
	CurrentCardNumber int
	// ProductList       []repository.Product
	ProductList []ProductInfo
}

type CardManager struct {
	Cards       map[string]*Card
	CurrentCard *Card
}

func NewCardManager() *CardManager {
	return &CardManager{
		Cards: make(map[string]*Card),
	}
}

func (cm *CardManager) GetCardByType(data string) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitProductService(db)

	products, err := service.GetProductByType(data)
	if err != nil {
		log.Println(err)
	}
	color.Redln(products)
	productInfos := make([]ProductInfo, len(products))
	for i, product := range products {
		productInfos[i] = &product
	}
	color.Redln(productInfos)
	cm.UpdateInfo(data, NewCard(productInfos))

	return nil
}

func (cm *CardManager) GetCartItemsByUserID(data string, userID int) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitCartService(db)

	products, err := service.GetItemsByUserID(int64(userID))
	if err != nil {
		log.Println(err)
	}

	productInfos := make([]ProductInfo, len(products))
	for i, product := range products {
		productInfos[i] = product
	}
	cm.UpdateInfo(data, NewCard(productInfos))

	return nil

}

func (cm *CardManager) GetCardAll(data string) error {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
	}
	defer db.Close()

	service := InitProductService(db)

	products, err := service.GetAllProducts()
	if err != nil {
		log.Println(err)
	}

	productInfos := make([]ProductInfo, len(products))
	for i, product := range products {
		productInfos[i] = &product
	}
	cm.UpdateInfo(data, NewCard(productInfos))

	return nil

}

func (cm *CardManager) NextCard() {
	cm.CurrentCard.next()
}

func (cm *CardManager) PrevCard() {
	cm.CurrentCard.prev()
}

func (cm *CardManager) GetCardText() string {
	return cm.CurrentCard.getTextTemplate()
}

func (cm *CardManager) GetCardImage() string {
	if cm.CurrentCard == nil {
		return ""
	}
	return cm.CurrentCard.Image
}

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

func NewCard(products []ProductInfo) *Card {
	color.Redln(products)
	if len(products) > 0 {
		return &Card{
			ID:                products[0].GetID(),
			Name:              products[0].GetName(),
			Type:              products[0].GetType(),
			Description:       products[0].GetDescription(),
			Quantity:          0,
			Price:             products[0].GetPrice(),
			Image:             products[0].GetImage(),
			TotalCards:        len(products),
			CurrentCardNumber: 0,
			ProductList:       products,
		}
	}
	return nil
}

func (c *Card) prev() {
	if c.CurrentCardNumber > 0 {
		c.CurrentCardNumber--
		c.ID = c.ProductList[c.CurrentCardNumber].GetID()
		c.Type = c.ProductList[c.CurrentCardNumber].GetType()
		c.Name = c.ProductList[c.CurrentCardNumber].GetName()
		c.Description = c.ProductList[c.CurrentCardNumber].GetDescription()
		c.Price = c.ProductList[c.CurrentCardNumber].GetPrice()
		c.Image = c.ProductList[c.CurrentCardNumber].GetImage()
	}

}
func (c *Card) next() {
	if c.CurrentCardNumber < c.TotalCards-1 {
		c.CurrentCardNumber++
		c.ID = c.ProductList[c.CurrentCardNumber].GetID()
		c.Type = c.ProductList[c.CurrentCardNumber].GetType()
		c.Name = c.ProductList[c.CurrentCardNumber].GetName()
		c.Description = c.ProductList[c.CurrentCardNumber].GetDescription()
		c.Price = c.ProductList[c.CurrentCardNumber].GetPrice()
		c.Image = c.ProductList[c.CurrentCardNumber].GetImage()
	}

}

func (c *Card) getTextTemplate() string {
	if c == nil {
		return "В корзине нет товаров"
	}
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCardNumber+1, c.TotalCards)
}
