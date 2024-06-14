package services

import (
	"fmt"
	"strings"

	repository "github.com/Resolution-hash/shop_bot/repository/cart"
)

type CartService struct {
	Repo repository.CartRepo
}

func NewCartService(repo repository.CartRepo) *CartService {
	return &CartService{
		repo,
	}
}

func (s *CartService) DeleteItem(item repository.CartItem) error {
	err := s.Repo.DeleteItem(item)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) GetQuantityByItemID(item repository.CartItem) (int, error) {
	quantity, err := s.Repo.GetQuantityByItemID(item)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (s *CartService) AddItem(item repository.CartItem) (int, error) {
	return s.Repo.AddItem(item)
}

func (s *CartService) Increment(item repository.CartItem) (int, error) {
	return s.Repo.Increment(item)
}

func (s *CartService) Decrement(item repository.CartItem) (int, error) {
	return s.Repo.Decrement(item)
}

func (s *CartService) GetItemsByUserID(userID int64) ([]*repository.CartProduct, error) {
	products, err := s.Repo.GetItemsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *CartService) FormatCartText(items []*repository.CartProduct) string {
	var sb strings.Builder
	totalPrice := 0.0

	for i, item := range items {
		itemTotal := item.Price * float64(item.Quantity)
		sb.WriteString(fmt.Sprintf("<b>%d</b>.%s\n   Количество: %d шт\n   Цена за 1 шт: %0.fруб.\n", i+1, item.Name, item.Quantity, item.Price))
		sb.WriteString("  ----------------------------------------------\n")
		totalPrice += itemTotal
	}

	sb.WriteString(fmt.Sprintf("\n💰 Итог: %0.fруб.", totalPrice))

	return sb.String()
}
