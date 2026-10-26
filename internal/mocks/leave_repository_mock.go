package mocks

import (
	"context"
	"hris-tools/internal/models"
	"hris-tools/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockLeaveRepository struct {
	mock.Mock
}

func (m *MockLeaveRepository) Create(leave *models.LeaveRequest) error {
	args := m.Called(leave)

	return args.Error(0)
}

func (m *MockLeaveRepository) LeaveByID(id uint) (*models.LeaveRequest, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LeaveRequest), args.Error(1)
}

func (m *MockLeaveRepository) LeaveByUserID(ctx context.Context, userID uint) ([]models.LeaveRequest, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.LeaveRequest), args.Error(1)
}

func (m *MockLeaveRepository) DeleteByID(id uint) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockLeaveRepository) LeaveTypeByID(id uint) (*models.LeaveType, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.LeaveType), args.Error(1)
}

func (m *MockLeaveRepository) LeaveBalanceByID(userID uint) (*models.LeaveBalance, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.LeaveBalance), args.Error(1)
}

func (m *MockLeaveRepository) GetAllPending(ctx context.Context) ([]models.LeaveRequest, error) {
	args := m.Called(ctx)

	return args.Get(0).([]models.LeaveRequest), args.Error(1)
}

func (m *MockLeaveRepository) ApprovedLeave(ctx context.Context, param *repository.ApproveLeaveParams) error {
	args := m.Called(ctx, param)

	return args.Error(0)
}

func (m *MockLeaveRepository) RejectedLeave(ctx context.Context, leaveID uint, hrID uint) error {
	args := m.Called(ctx, leaveID, hrID)

	return args.Error(0)
}
