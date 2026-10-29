package service_test

import (
	"hris-tools/internal/dto"
	"hris-tools/internal/mocks"
	"hris-tools/internal/service"
	"testing"

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
