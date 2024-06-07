package services

import repository "github.com/Resolution-hash/shop_bot/internal/repository/admin"

type AdminService struct {
	repo repository.AdminRepo
}

func NewAdminService(repo repository.AdminRepo) *AdminService {
	return &AdminService{
		repo,
	}
}

// AddAdmin(int64) error
// 	deleteAdmin(int64) error
// 	IsAdmin(int64) bool

func (s *AdminService) AddAdmin(userID int64) error {
	err := s.repo.AddAdmin(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) DeleteAdmin(userID int64) error {
	err := s.repo.DeleteAdmin(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) IsAdmin(userID int64) (bool, error) {
	isAdmin, err := s.repo.IsAdmin(userID)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}
