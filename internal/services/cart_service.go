package services

import (
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
