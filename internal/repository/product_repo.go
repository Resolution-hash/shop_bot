package repository

import "reflect"

type Product struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Price       float64
	Image string
}

func IsEmpty(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

type ProductRepo interface {
	CreateProduct(Product) error
	UpdateProduct(Product) error
	DeleteProduct(int64) error
	GetAllProducts() ([]Product, error)
	GetProductsByType(productType string) ([]Product, error)
}
