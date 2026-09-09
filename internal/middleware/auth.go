package middleware

import (
	"errors"
	"hris-tools/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ") // melepas token dari bearer

		if len(tokenString) != 2 || tokenString[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			c.Abort()
			return
		}

		token := tokenString[1]
		// validate token
		claims, err := utils.ValidateToken(token)

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
				})
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Role)
		c.Next()
	}
}
