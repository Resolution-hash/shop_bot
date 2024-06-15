package services

import (
	repository "github.com/Resolution-hash/shop_bot/repository/product"
)

type ProductService struct {
	repo repository.ProductRepo
}

func NewProductService(repo repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product repository.Product) error {
	err := s.repo.CreateProduct(product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(ID int64) error {
	err := s.repo.DeleteProduct(ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) UpdateProductImage(product repository.Product) error {
	err := s.repo.UpdateProductImage(product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) UpdateProductText(product repository.Product) error {
	err := s.repo.UpdateProductText(product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetAllProducts() ([]repository.Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByType(productType string) ([]repository.Product, error) {
	products, err := s.repo.GetProductsByType(productType)
	if err != nil {
		return nil, err
	}
	return products, nil
}
