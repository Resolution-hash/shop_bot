package services

import "github.com/Resolution-hash/shop_bot/internal/repository/user"

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo,
	}
}

func (u *UserService) AddUser(user repository.User) error {
	err := u.repo.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}
