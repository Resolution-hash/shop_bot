package repository

type CartItem struct {
	ProductID int64
	UserID    int64
	Quantity  int
}

type CartProduct struct {
	ProductID   int
	Name        string
	Type        string
	Description string
	Price       float64
	Quantity    int
	Image       string
}

func (c *CartProduct) GetID() int64 {
	return int64(c.ProductID)
}
func (c *CartProduct) GetName() string {
	return c.Name
}
func (c *CartProduct) GetType() string {
	return c.Type
}
func (c *CartProduct) GetDescription() string {
	return c.Description
}
func (c *CartProduct) GetPrice() float64 {
	return c.Price
}
func (c *CartProduct) GetImage() string {
	return c.Image
}

type CartRepo interface {
	AddItem(CartItem) (int, error)
	Increment(CartItem) (int, error)
	Decrement(CartItem) (int, error)
	GetQuantityByItemID(CartItem) (int, error)
	GetItemsByUserID(int64) ([]*CartProduct, error)
}
