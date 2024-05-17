package services

import (
	"log"

	"github.com/Resolution-hash/shop_bot/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepo
}

func NewProductService(repo repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]repository.Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		log.Println("Error getting all products: ", err)
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByType(productType string) ([]repository.Product, error) {
	products, err := s.repo.GetProductsByType(productType)
	if err != nil {
		log.Println("Error getting all products: ", err)
		return nil, err
	}
	return products, nil
}
