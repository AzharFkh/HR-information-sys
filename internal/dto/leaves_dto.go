package dto

import "hris-tools/internal/models"

type CreateLeaveRequestReq struct {
	UserID      uint   `json:"user_id" binding:"required"`
	LeaveTypeID uint   `json:"leave_type_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // Terima sebagai string (YYYY-MM-DD)
	EndDate     string `json:"end_date" binding:"required"`   // Terima sebagai string (YYYY-MM-DD)
	Reason      string `json:"reason" binding:"required,max=200"`
}

type LeaveResponse struct {
	ID           uint               `json:"id"`
	UserID       uint               `json:"user_id"`
	LeaveTypeID  uint               `json:"leave_type_id"`
	StartDate    string             `json:"start_date"`
	EndDate      string             `json:"end_date"`
	Reason       string             `json:"reason"`
	Status       models.LeaveStatus `json:"status"`
	ApprovedByID *uint              `json:"approved_by_id"`

	User      UserResponse      `json:"user"`
	LeaveType LeaveTypeResponse `json:"leave_type"`
}

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LeaveTypeResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	IsDeductQuota bool   `json:"is_deduct_quota"`
}

func ToLeaveResponse(leave models.LeaveRequest) LeaveResponse {
	return LeaveResponse{
		ID:           leave.ID,
		UserID:       leave.UserID,
		LeaveTypeID:  leave.LeaveTypeID,
		StartDate:    leave.StartDate.Format("2006-12-31"), // (yyyy-mm-dd)
		EndDate:      leave.EndDate.Format("2006-12-31"),
		Reason:       leave.Reason,
		Status:       leave.Status,
		ApprovedByID: leave.ApprovedByID,

		User: UserResponse{
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
