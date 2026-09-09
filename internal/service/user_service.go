package service

import (
	"hris-tools/internal/models"
	"hris-tools/internal/repository"
)

type UserService interface {
	CreateUser(user models.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// wasnt finish
func (s *userService) CreateUser(user models.User) error {
	return nil
}

// user service

