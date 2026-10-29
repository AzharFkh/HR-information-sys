package dto

import "hris-tools/internal/models"

type CreateLeaveRequestReq struct {
	UserID      uint   `json:"-"`
	LeaveTypeID uint   `json:"leave_type_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // Terima sebagai string (YYYY-MM-DD)
	EndDate     string `json:"end_date" binding:"required"`   // Terima sebagai string (YYYY-MM-DD)
	Reason      string `json:"reason" binding:"required,max=200"`
}

// for HR
type LeaveResponse struct {
	ID           uint               `json:"id"`
	StartDate    string             `json:"start_date"`
	EndDate      string             `json:"end_date"`
	Reason       string             `json:"reason"`
	Status       models.LeaveStatus `json:"status"`
	ApprovedByID *uint              `json:"approved_by_id"`

	User      UserInformation   `json:"user"`
	LeaveType LeaveTypeResponse `json:"leave_type"`
}

type LeaveTypeResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	IsDeductQuota bool   `json:"is_deduct_quota"`
}

func ToLeaveResponse(leave models.LeaveRequest) LeaveResponse {
	return LeaveResponse{
		ID:           leave.ID,
		StartDate:    leave.StartDate.Format("2006-01-02"), // (yyyy-mm-dd)
		EndDate:      leave.EndDate.Format("2006-01-02"),
		Reason:       leave.Reason,
		Status:       leave.Status,
		ApprovedByID: leave.ApprovedByID,

		User: UserInformation{
			ID:   leave.UserID,
			Name: leave.User.Name,
		},

		LeaveType: LeaveTypeResponse{
			ID:            leave.LeaveTypeID,
			Name:          leave.LeaveType.Name,
			IsDeductQuota: leave.LeaveType.IsDeductQuota,
		},
	}
}

// untuk return banyak response
func ToLeaveResponseList(leaves []models.LeaveRequest) []LeaveResponse {
	var responses []LeaveResponse
	for _, leave := range leaves {
		responses = append(responses, ToLeaveResponse(leave))
	}
	return responses
}

// for user leave request
type CreateLeaveResponse struct {
	ID          uint   `json:"id"`
	LeaveTypeID uint   `json:"leave_type_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Reason      string `json:"reason"`
	Status      string `json:"status"`
}

// Fungsi Mapper untuk menyaring Model menjadi DTO
func ToCreateLeaveResponse(leave *models.LeaveRequest) CreateLeaveResponse {
	return CreateLeaveResponse{
		ID:          leave.ID,
		LeaveTypeID: leave.LeaveTypeID,
		// Ubah format waktu kembali menjadi YYYY-MM-DD agar cantik di JSON
		StartDate: leave.StartDate.Format("2006-01-02"),
		EndDate:   leave.EndDate.Format("2006-01-02"),
		Reason:    leave.Reason,
		Status:    string(leave.Status), // Ubah custom type LeaveStatus menjadi string
	}
}

type UpdateLeaveStatusReq struct {
	HRID    uint               `json:"-"`
	LeaveID uint               `json:"-"`
	Status  models.LeaveStatus `json:"status" binding:"required,oneof=APPROVED REJECTED"`
}
