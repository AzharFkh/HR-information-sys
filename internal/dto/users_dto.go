package dto

import (
	"hris-tools/internal/models"
	"time"
)

type UserRegister struct {
	Name     string
	Password string
}

type UserInformation struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID        uint      `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(user models.User) UserResponse {

	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
