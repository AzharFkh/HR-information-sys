package service

import (
	"context"
	"errors"
	"fmt"
	"hris-tools/internal/dto"
	"hris-tools/internal/models"
	"hris-tools/internal/repository"
	"time"
)

type LeaveService interface {
	Create(leave dto.CreateLeaveRequestReq) (*models.LeaveRequest, error)
	GetLeavesByPending(ctx context.Context) ([]models.LeaveRequest, error)
	GetLeavesByUser(ctx context.Context, userID uint) ([]models.LeaveRequest, error)
	UpdateStatus(ctx context.Context, req *dto.UpdateLeaveStatusReq) error
}

type leaveService struct {
	repo repository.LeavesRepository
}

func NewLeaveService(repo repository.LeavesRepository) LeaveService {
	return &leaveService{repo}
}

func (s *leaveService) Create(req dto.CreateLeaveRequestReq) (*models.LeaveRequest, error) {

	dateLayout := "2006-01-02"

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, errors.New("format tanggal mulai salah, gunakan YYYY-MM-DD")
	}

	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, errors.New("format tanggal selesai salah, gunakan YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, errors.New("Tanggal mulai tidak boleh lebih dari tanggal selesai")
	}

	durationInHours := endDate.Sub(startDate).Hours()
	requestedDays := int(durationInHours/24) + 1

	leaveType, err := s.repo.LeaveTypeByID(req.LeaveTypeID)

	if err != nil {
		return nil, errors.New("Jenis cuti tidak ditemukan")
	}

	if leaveType.IsDeductQuota {
		leaveBalance, err := s.repo.LeaveBalanceByID(req.UserID)
		if err != nil {
			return nil, errors.New("jatah cuti user tidak dapat ditemukan")
		}

		if leaveBalance.Quota < requestedDays {
			return nil, errors.New("jatah cuti tidak mencukupi")
		}
	}

	newLeave := models.LeaveRequest{
		UserID:      req.UserID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   startDate,
		EndDate:     endDate,
		Reason:      req.Reason,
		Status:      models.StatusPending,
		LeaveType:   *leaveType,
	}

	err = s.repo.Create(&newLeave)
	if err != nil {
		return nil, errors.New("gagal menyimpan pengajuan ke database")
	}

	return &newLeave, nil
}

func (s *leaveService) GetLeavesByPending(ctx context.Context) ([]models.LeaveRequest, error) {
	leaves, err := s.repo.GetAllPending(ctx)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pending: %w", err) //gunakan %w
	}

	return leaves, nil // slice models
}

func (s *leaveService) GetLeavesByUser(ctx context.Context, userID uint) ([]models.LeaveRequest, error) {
	leaves, err := s.repo.LeaveByUserID(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pending: %w", err)
	}
	return leaves, nil
}

func (s *leaveService) UpdateStatus(ctx context.Context, req *dto.UpdateLeaveStatusReq) error {

	leave, err := s.repo.LeaveByID(req.LeaveID)

	if err != nil {
		return errors.New("pengajuan cuti tidak ditemukan")
	}

	if leave.Status != models.StatusPending {
		return errors.New("hanya pengajuan berstatus PENDING yang dapat diproses")
	}

	durationInHours := leave.EndDate.Sub(leave.StartDate).Hours()
	requestedDays := uint(durationInHours/24) + 1 // +1 karena hari hari dimulai

	switch req.Status {
	case models.StatusApproved:
		// return s.repo.ApprovedLeave(ctx, req.LeaveID, leave.UserID, requestedDays)
		params := repository.ApproveLeaveParams{
			LeaveID: req.LeaveID,
			HRID: req.HRID,
			EmployeeID: leave.UserID,
			ReqDays: requestedDays,
		}
		
		return s.repo.ApprovedLeave(ctx, &params)
	case models.StatusRejected:
		return s.repo.RejectedLeave(ctx, req.LeaveID, req.HRID)
	default:
		return errors.New("status tidak valid")
	}
}
