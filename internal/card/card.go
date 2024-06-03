package card

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/gookit/color"
)

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
	ProductList       []repository.Product
}

type CardManager struct {
	Cards              map[string]*Card
	CardsProductCart   map[string]*CardProductCart
	CurrentCard        *Card
	CurrentProductCart *CardProductCart
}

func NewCardManager() *CardManager {
	return &CardManager{
		Cards:            make(map[string]*Card),
		CardsProductCart: make(map[string]*CardProductCart),
	}
}

func (cm *CardManager) GetCardOnCart(userID int) error {
	db, err := setupDatabase()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	service := initCardService(db)

	products, err := service.GetItemsByUserID(int64(userID))
	if err != nil {
		log.Println(err)
	}
	cm.UpdateInfo(data, NewCard(products))

	return nil
}

func (cm *CardManager) GetCardByType(data string) error {
	db, err := setupDatabase()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	service := initProductService(db)

	products, err := service.GetProductByType(data)
	if err != nil {
		log.Println(err)
	}
	cm.UpdateInfo(data, NewCard(products))

	return nil

	// if repository.IsEmpty(products) {
	// 	return nil, fmt.Errorf("")
	// 	// keyboard := messages.GetKeyboard("back", "Магазин")
	// 	// messages.SendMessage(bot, userInfo.UserID, "Товаров нет.", keyboard)
	// 	// return
	// }

}

func (cm *CardManager) GetCardAll(data string) error {
	db, err := setupDatabase()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	service := initProductService(db)

	products, err := service.GetAllProducts()
	if err != nil {
		log.Println(err)
	}
	cm.UpdateInfo(data, NewCard(products))

	return nil

	// if repository.IsEmpty(products) {
	// 	return nil, fmt.Errorf("")
	// 	// keyboard := messages.GetKeyboard("back", "Магазин")
	// 	// messages.SendMessage(bot, userInfo.UserID, "Товаров нет.", keyboard)
	// 	// return
	// }

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
	return cm.CurrentCard.Image
}

func (cm *CardManager) UpdateInfo(step string, card *Card) {
	cm.Cards[step] = card
	cm.CurrentCard = card
}

func (cm *CardManager) UpdateInfoCart(step string, card *CardProductCart) {
	cm.CardsProductCart[step] = card
	cm.CurrentProductCart = card
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
			ID:                products[0].ID,
			Name:              products[0].Name,
			Type:              products[0].Type,
			Description:       products[0].Description,
			Quantity:          0,
			Price:             products[0].Price,
			Image:             products[0].Image,
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
		c.ID = c.ProductList[c.CurrentCardNumber].ID
		c.Type = c.ProductList[c.CurrentCardNumber].Type
		c.Name = c.ProductList[c.CurrentCardNumber].Name
		c.Description = c.ProductList[c.CurrentCardNumber].Description
		c.Price = c.ProductList[c.CurrentCardNumber].Price
		c.Image = c.ProductList[c.CurrentCardNumber].Image
	}

}
func (c *Card) next() {
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

func (c *Card) getTextTemplate() string {
	return fmt.Sprintf("Название: %s\n\nОписание: %s\n\nЦена: %0.f рублей\n%d/%d", c.Name, c.Description, c.Price, c.CurrentCardNumber+1, c.TotalCards)
}

func initProductService(db *sql.DB) *services.ProductService {
	repo := repository.NewSqliteProductRepo(db)
	service := services.NewProductService(repo)
	return service
}

func setupDatabase() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		fmt.Println("error to get cfg.DbUrl")
		return nil, err
	}
	return db, nil
}
