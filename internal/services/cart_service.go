package services

import (
	"fmt"
	"strings"

	"github.com/Resolution-hash/shop_bot/internal/repository"
)

type CartService struct {
	Repo repository.CartRepo
}

func NewCartService(repo repository.CartRepo) *CartService {
	return &CartService{
		repo,
	}
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

func (s *CartService) GetCartText(userID int64) (string, error) {
	products, err := s.GetItemsByUserID(userID)
	if err != nil {
		return "", err
	}

	return formatCartText(products), nil
}

func (s *CartService) GetItemsByUserID(userID int64) ([]*repository.CartProduct, error) {
	products, err := s.Repo.GetItemsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func formatCartText(items []*repository.CartProduct) string {
	var sb strings.Builder
	totalPrice := 0.0

	for i, item := range items {
		itemTotal := item.Price * float64(item.Quantity)
		sb.WriteString(fmt.Sprintf("<b>%d</b>.%s\n   Количество: %dшт\n   Цена за 1 шт: <i>%0.fруб</i>.\n", i+1, item.Name, item.Quantity, item.Price))
		sb.WriteString(fmt.Sprint("  ----------------------\n"))
		totalPrice += itemTotal
	}

	sb.WriteString(fmt.Sprintf("\n💰 Итог: %0.fруб.", totalPrice))

	return sb.String()
}
