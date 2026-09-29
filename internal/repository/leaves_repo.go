package repository

import (
	"context"
	"hris-tools/internal/models"

	"gorm.io/gorm"
)

// bisa ditambahkan ctx context.Context
type LeavesRepository interface {
	Create(leave *models.LeaveRequest) error
	LeaveByID(id uint) (*models.LeaveRequest, error)
	LeaveByUserID(ctx context.Context, userID uint) ([]models.LeaveRequest, error)
	DeleteByID(id uint) error
	LeaveTypeByID(id uint) (*models.LeaveType, error)
	LeaveBalanceByID(userID uint) (*models.LeaveBalance, error)
	GetAllPending(ctx context.Context) ([]models.LeaveRequest, error)
	ApprovedLeave(ctx context.Context, param *ApproveLeaveParams) error
	RejectedLeave(ctx context.Context, leaveID uint, hrID uint) error
}

type leavesRepository struct {
	db *gorm.DB
}

func NewLeavesRepository(db *gorm.DB) LeavesRepository {
	return &leavesRepository{db}
}

func (r *leavesRepository) Create(leave *models.LeaveRequest) error {
	return r.db.Create(leave).Error
}

func (r *leavesRepository) LeaveByID(id uint) (*models.LeaveRequest, error) {
	var leave models.LeaveRequest
	err := r.db.First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *leavesRepository) LeaveByUserID(ctx context.Context, userID uint) ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("LeaveType").
		Where("user_id = ?", userID).
		Find(&leaves).Error

	if err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *leavesRepository) DeleteByID(id uint) error {
	var leave models.LeaveRequest

	return r.db.Delete(leave).Error
}

// pr
func (r *leavesRepository) LeaveTypeByID(id uint) (*models.LeaveType, error) {
	var leaveType models.LeaveType

	err := r.db.First(&leaveType, id).Error
	if err != nil {
		return nil, err
	}

	return &leaveType, nil

}

func (r *leavesRepository) LeaveBalanceByID(userID uint) (*models.LeaveBalance, error) {
	// wadah kosong
	var leave models.LeaveBalance
	// querry
	err := r.db.Where("user_id = ?", userID).First(&leave).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *leavesRepository) GetAllPending(ctx context.Context) ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest

	err := r.db.WithContext(ctx). // Withcontext
					Preload("User").
					Preload("LeaveType").
					Where("status = ?", models.StatusPending).
					Find(&leaves).Error

	if err != nil {
		return nil, err
	}
	return leaves, nil
}

type ApproveLeaveParams struct {
	LeaveID    uint
	HRID       uint
	EmployeeID uint
	ReqDays    uint
}

func (r *leavesRepository) ApprovedLeave(ctx context.Context, param *ApproveLeaveParams) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// var leaves models.LeaveRequest

		if err := tx.Model(&models.LeaveRequest{}).
			Where("id = ?", param.LeaveID).
			Updates(map[string]any{
				"status":         models.StatusApproved,
				"approved_by_id": param.HRID,
			}).Error; err != nil {
			return err
		}

		// insert here
		if err := tx.Model(&models.LeaveBalance{}).
			Where("user_id = ?", param.EmployeeID).
			Update("used", gorm.Expr("used + ?", param.ReqDays)).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *leavesRepository) RejectedLeave(ctx context.Context, leaveID uint, hrID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.LeaveRequest{}).
		Where("id = ?", leaveID).
		Updates(map[string]any{
			"status":         models.StatusRejected,
			"approved_by_id": hrID,
		}).Error
}
