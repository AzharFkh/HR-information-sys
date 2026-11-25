package service_test

import (
	"context"
	"hris-tools/internal/dto"
	"hris-tools/internal/mocks"
	"hris-tools/internal/models"
	"hris-tools/internal/repository"
	"hris-tools/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaveService_Create_ValidationDates(t *testing.T) {

	tests := []struct {
		name        string
		request     dto.CreateLeaveRequestReq
		expectedErr string
	}{
		{
			name: "Error - Format StartDate salah (menggunakan DD-MM-YYYY)",
			request: dto.CreateLeaveRequestReq{
				StartDate: "25-12-2023",
				EndDate:   "2023-12-30",
			},
			expectedErr: "format tanggal mulai salah, gunakan YYYY-MM-DD",
		},
		{
			name: "Error - Format EndDate salah (menggunakan slash)",
			request: dto.CreateLeaveRequestReq{
				StartDate: "2023-12-25",
				EndDate:   "2023/12/30",
			},
			expectedErr: "format tanggal selesai salah, gunakan YYYY-MM-DD",
		},

		{
			name: "Error - Logika terbalik (StartDate lebih besar dari EndDate)",
			request: dto.CreateLeaveRequestReq{
				StartDate: "2023-12-30",
				EndDate:   "2023-12-25",
			},
			expectedErr: "Tanggal mulai tidak boleh lebih dari tanggal selesai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaveRepo := new(mocks.MockLeaveRepository)
			leaveService := service.NewLeaveService(leaveRepo)
			result, err := leaveService.Create(tt.request)

			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
			assert.Nil(t, result)

			leaveRepo.AssertNotCalled(t, "LeaveTypeByID", mock.Anything)
			leaveRepo.AssertNotCalled(t, "Create", mock.Anything)
		})
	}
}

func TestLeaveService_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2026, 11, 23, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC) // 3 hari

	t.Run("Approve cuti berkuota ketika saldo mencukupi", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:        1,
			UserID:    10,
			StartDate: start,
			EndDate:   end,
			Status:    models.StatusPending,
			LeaveType: models.LeaveType{ID: 1, Name: "Cuti Tahunan", IsDeductQuota: true},
		}
		balance := &models.LeaveBalance{UserID: 10, Quota: 12, Used: 5}

		leaveRepo.On("LeaveByID", uint(1)).Return(leave, nil)
		leaveRepo.On("LeaveBalanceByID", uint(10)).Return(balance, nil)
		leaveRepo.On("ApprovedLeave", ctx, &repository.ApproveLeaveParams{
			LeaveID:     1,
			HRID:        99,
			EmployeeID:  10,
			ReqDays:     3,
			DeductQuota: true,
		}).Return(nil)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 1,
			Status:  models.StatusApproved,
		})

		assert.NoError(t, err)
		leaveRepo.AssertExpectations(t)
	})

	t.Run("Approve cuti non-kuota tetap berhasil tanpa cek atau potong saldo", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:        2,
			UserID:    10,
			StartDate: start,
			EndDate:   end,
			Status:    models.StatusPending,
			LeaveType: models.LeaveType{ID: 2, Name: "Sakit", IsDeductQuota: false},
		}

		leaveRepo.On("LeaveByID", uint(2)).Return(leave, nil)
		leaveRepo.On("ApprovedLeave", ctx, &repository.ApproveLeaveParams{
			LeaveID:     2,
			HRID:        99,
			EmployeeID:  10,
			ReqDays:     3,
			DeductQuota: false,
		}).Return(nil)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 2,
			Status:  models.StatusApproved,
		})

		assert.NoError(t, err)
		leaveRepo.AssertNotCalled(t, "LeaveBalanceByID", mock.Anything)
		leaveRepo.AssertNotCalled(t, "RejectedLeave", mock.Anything, mock.Anything, mock.Anything)
		leaveRepo.AssertExpectations(t)
	})

	t.Run("Approve cuti berkuota ditolak otomatis jika kuota tidak cukup", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:        3,
			UserID:    10,
			StartDate: start,
			EndDate:   end,
			Status:    models.StatusPending,
			LeaveType: models.LeaveType{ID: 1, Name: "Cuti Tahunan", IsDeductQuota: true},
		}
		balance := &models.LeaveBalance{UserID: 10, Quota: 12, Used: 10}

		leaveRepo.On("LeaveByID", uint(3)).Return(leave, nil)
		leaveRepo.On("LeaveBalanceByID", uint(10)).Return(balance, nil)
		leaveRepo.On("RejectedLeave", ctx, uint(3), uint(99)).Return(nil)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 3,
			Status:  models.StatusApproved,
		})

		assert.Error(t, err)
		assert.Equal(t, "gagal disetujui: kuota cuti karyawan tidak mencukupi", err.Error())
		leaveRepo.AssertNotCalled(t, "ApprovedLeave", mock.Anything, mock.Anything)
		leaveRepo.AssertExpectations(t)
	})

	t.Run("Reject cuti tidak perlu cek saldo kuota", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:     4,
			UserID: 10,
			Status: models.StatusPending,
		}

		leaveRepo.On("LeaveByID", uint(4)).Return(leave, nil)
		leaveRepo.On("RejectedLeave", ctx, uint(4), uint(99)).Return(nil)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 4,
			Status:  models.StatusRejected,
		})

		assert.NoError(t, err)
		leaveRepo.AssertNotCalled(t, "LeaveBalanceByID", mock.Anything)
		leaveRepo.AssertExpectations(t)
	})

	t.Run("Race condition: approve gagal jika status sudah tidak PENDING di transaksi", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:        5,
			UserID:    10,
			StartDate: start,
			EndDate:   end,
			Status:    models.StatusPending,
			LeaveType: models.LeaveType{ID: 1, Name: "Cuti Tahunan", IsDeductQuota: true},
		}
		balance := &models.LeaveBalance{UserID: 10, Quota: 12, Used: 5}

		leaveRepo.On("LeaveByID", uint(5)).Return(leave, nil)
		leaveRepo.On("LeaveBalanceByID", uint(10)).Return(balance, nil)
		leaveRepo.On("ApprovedLeave", ctx, &repository.ApproveLeaveParams{
			LeaveID:     5,
			HRID:        99,
			EmployeeID:  10,
			ReqDays:     3,
			DeductQuota: true,
		}).Return(repository.ErrLeaveNotPending)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 5,
			Status:  models.StatusApproved,
		})

		assert.Error(t, err)
		assert.Equal(t, "hanya pengajuan berstatus PENDING yang dapat diproses", err.Error())
		leaveRepo.AssertExpectations(t)
	})

	t.Run("Race condition: approve gagal jika kuota sudah terpakai transaksi lain", func(t *testing.T) {
		leaveRepo := new(mocks.MockLeaveRepository)
		leaveService := service.NewLeaveService(leaveRepo)

		leave := &models.LeaveRequest{
			ID:        6,
			UserID:    10,
			StartDate: start,
			EndDate:   end,
			Status:    models.StatusPending,
			LeaveType: models.LeaveType{ID: 1, Name: "Cuti Tahunan", IsDeductQuota: true},
		}
		balance := &models.LeaveBalance{UserID: 10, Quota: 12, Used: 5}

		leaveRepo.On("LeaveByID", uint(6)).Return(leave, nil)
		leaveRepo.On("LeaveBalanceByID", uint(10)).Return(balance, nil)
		leaveRepo.On("ApprovedLeave", ctx, &repository.ApproveLeaveParams{
			LeaveID:     6,
			HRID:        99,
			EmployeeID:  10,
			ReqDays:     3,
			DeductQuota: true,
		}).Return(repository.ErrInsufficientQuota)

		err := leaveService.UpdateStatus(ctx, &dto.UpdateLeaveStatusReq{
			HRID:    99,
			LeaveID: 6,
			Status:  models.StatusApproved,
		})

		assert.Error(t, err)
		assert.Equal(t, "gagal disetujui: kuota cuti karyawan tidak mencukupi", err.Error())
		leaveRepo.AssertExpectations(t)
	})
}
