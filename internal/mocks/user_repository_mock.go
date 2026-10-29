package mocks

import (
	"hris-tools/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(name string) (*models.Role, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

// Mock repo for user

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User, balance *models.LeaveBalance) error {
	args := m.Called(user, balance)

	return args.Error(0)
}

func (m *MockUserRepository) UserByID(id uint) (*models.User, error) {
	args := m.Called(id)

	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UserByName(name string) (*models.User, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}
