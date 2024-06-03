package repository

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
	Image       string
}

type CartRepo interface {
	AddItem(CartItem) (int, error)
	Increment(CartItem) (int, error)
	Decrement(CartItem) (int, error)
	GetQuantityByItemID(CartItem) (int, error)
	// RemoveItem(int64) error
	GetItemsByUserID(int64) ([]*CartProduct, error)
}
