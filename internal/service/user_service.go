package service

import (
	"errors"
	"hris-tools/internal/dto"
	"hris-tools/internal/models"
	"hris-tools/internal/repository"
	"hris-tools/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req dto.UserRegister) (*models.User, error)
	Login(req dto.UserRegister) (string, error)
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
		repo:     repo,
		roleRepo: roleRepo,
	}
}

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
		Name:     req.Name,
		Password: string(hashed),
		RoleID:   role.ID,
	}

	leaveBalance := models.LeaveBalance{Quota: 12} // standard 12

	// insert to db
	if err := s.repo.Create(&user, &leaveBalance); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) Login(req dto.UserRegister) (string, error) {

	user, err := s.repo.UserByName(req.Name)

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(utils.TokenPayload{
		UserID: user.ID,
		Role:   user.Role.Name,
	})

	if err != nil {
		return "", err
	}
	return token, nil
}
