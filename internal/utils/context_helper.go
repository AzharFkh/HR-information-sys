package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, error) {
	userIDCtx, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("unauthorized")
	}

	switch v := userIDCtx.(type) {
	case float64:
		return uint(v), nil
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	default:
		return 0, errors.New("invalid user id type in token")
	}
}
