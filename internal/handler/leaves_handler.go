package handler

import (
	"context"
	"errors"
	"fmt"
	"hris-tools/internal/dto"
	"hris-tools/internal/service"
	"hris-tools/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	leaveService service.LeaveService
}

func NewLeaveHandler(leaveService service.LeaveService) *LeaveHandler {
	return &LeaveHandler{leaveService}
}

// bisa tambahkan ctx context.Context (withtimeout)

func (h *LeaveHandler) CreateRequest(c *gin.Context) {
	var req dto.CreateLeaveRequestReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID

	leave, err := h.leaveService.Create(req)

	if err != nil {
		var clientErr *utils.ClientError

		if errors.As(err, &clientErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": clientErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "terjadi kesalahan saat membuat pengajuan cuti",
		})
		return
	}

	responseDTO := dto.ToCreateLeaveResponse(leave)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "leave request created",
		"data":    responseDTO,
	})

}

func (h *LeaveHandler) GetPendingRequest(c *gin.Context) {

	parentCtx := c.Request.Context()

	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)

	defer cancel()

	leaves, err := h.leaveService.GetLeavesByPending(ctx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"status":  "error",
				"message": "system took too long to response",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	leaveResponses := dto.ToLeaveResponseList(leaves)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "berhasil mengambil daftar antrean cuti",
		"data":    leaveResponses,
	})
}

func (h *LeaveHandler) GetAllLeave(c *gin.Context) {

	parentCtx := c.Request.Context()

	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)

	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	leaves, err := h.leaveService.GetLeavesByUser(ctx, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	leaveResponses := dto.ToLeaveResponseList(leaves)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "berhasil mengambil daftar pengajuan cuti",
		"data":    leaveResponses,
	})
}

func (h *LeaveHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	HRID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	idParam := c.Param("id")
	leaveID, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pengajuan tidak valid",
		})
		return
	}

	var req dto.UpdateLeaveStatusReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.HRID = HRID
	req.LeaveID = uint(leaveID)

	err = h.leaveService.UpdateStatus(ctx, &req)

	if err != nil {
		var clientErr *utils.ClientError

		if errors.As(err, &clientErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": clientErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Terjadi kesalahan pada server saat memproses pengajuan cuti",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("status pengajuan ID: %d berhasil diperbarui", req.LeaveID),
	})

}
