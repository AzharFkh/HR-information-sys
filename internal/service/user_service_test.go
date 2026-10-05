package service_test

import (
	"hris-tools/internal/dto"
	"hris-tools/internal/mocks"
	"hris-tools/internal/models"
	"hris-tools/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	req := dto.UserRegister{
		Name:     "Vicky",
		Password: "password1234",
	}

	mockRole := &models.Role{
		ID:   3,
		Name: "employee",
	}

	roleRepo.On("FindByName", "employee").Return(mockRole, nil)

	userRepo.On("Create", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.LeaveBalance")).Return(nil)

	result, err := userService.CreateUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Vicky", result.Name)
	assert.Equal(t, mockRole.ID, result.RoleID)

	assert.NotEqual(t, "password1234", result.Password)

	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	const username = "bambang"
	const plainPassword = "password4321"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	mockUser := &models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     username,
		Password: string(hashedPassword),
		RoleID:   3,

		Role: models.Role{
			ID:   3,
			Name: "employee",
		},
	}

	req := dto.UserRegister{
		Name:     username,
		Password: plainPassword,
	}

	userRepo.On("UserByName", username).Return(mockUser, nil)

	token, err := userService.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
}

func TestLogin_WrongPass(t *testing.T) {

	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)

	userService := service.NewUserService(userRepo, roleRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1234"), bcrypt.DefaultCost)

	mockUser := &models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "Bambang",
		Password: string(hashedPassword),
		RoleID:   3,

		Role: models.Role{
			ID:   3,
			Name: "employee",
		},
	}

	req := dto.UserRegister{Name: "bambang", Password: "salah-pass"}

	userRepo.On("UserByName", "bambang").Return(mockUser, nil)

	token, err := userService.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	userRepo.AssertExpectations(t)

}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	req := dto.UserRegister{Name: "tidak-ada", Password: "apa-saja"}

	userRepo.On("UserByName", "tidak-ada").Return(nil, gorm.ErrRecordNotFound)

	token, err := userService.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	userRepo.AssertExpectations(t)
}
