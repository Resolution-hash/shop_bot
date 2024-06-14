package repository

import (
	"reflect"
)

type Product struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Price       float64
	Image       string
}

func (p *Product) GetID() int64 {
	return p.ID
}
func (p *Product) GetName() string {
	return p.Name
}
func (p *Product) GetType() string {
	return p.Type
}
func (p *Product) GetDescription() string {
	return p.Description
}
func (p *Product) GetPrice() float64 {
	return p.Price
}
func (p *Product) GetImage() string {
	return p.Image
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
