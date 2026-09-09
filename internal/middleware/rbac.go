package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequiredRold(role string) gin.HandlerFunc{
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if userRole == role{
			return 
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})

	}
}