package service

import (
	"hris-tools/internal/dto"
	"hris-tools/internal/models"
	"hris-tools/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req dto.UserRegister) (*models.User, error)
}

type userService struct {
	repo     repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(
	repo repository.UserRepository,
	roleRepo repository.RoleRepository,
	) UserService {
	return &userService{
		repo: repo,
		roleRepo: roleRepo,
	}
}

// wasnt finish
func (s *userService) CreateUser(req dto.UserRegister) (*models.User, error) {
	// password hashing
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// insert default role
	role, err := s.roleRepo.FindByName("employee") // default as employee

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name: req.Name,
		Password: string(hashed),
		RoleID: role.ID,
	}

	// insert to db
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// user service
