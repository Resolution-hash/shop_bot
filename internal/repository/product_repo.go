package repository

type Product struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Price       float64
}

type ProductRepo interface {
	CreateProduct(Product) error
	UpdateProduct(Product) error
	DeleteProduct(int64) error
	GetAllProducts() ([]Product, error)
	GetProductsByType(productType string) ([]Product, error)
}
