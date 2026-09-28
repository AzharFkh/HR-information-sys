package handler

import (
	"context"
	"errors"
	"hris-tools/internal/dto"
	"hris-tools/internal/service"
	"net/http"
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

	userIDCtx, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint
    switch v := userIDCtx.(type) {
    case float64:
        userID = uint(v) // ubah ke float64
    case uint:
        userID = v       // gunakan bawaan uint
    case int:
        userID = uint(v)
    }

	req.UserID = userID

	leave, err := h.leaveService.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "leave request created",
		"data":    leave,
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
