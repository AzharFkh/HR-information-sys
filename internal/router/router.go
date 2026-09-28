package router

import (
	"hris-tools/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setuprouter(
	userHandler *handler.UserHandler,
	leaveHandler *handler.LeaveHandler,
) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "uji",
		})
	})

	r.POST("/leaves", leaveHandler.CreateRequest)
	r.GET("/leaves/pending", leaveHandler.GetPendingRequest)

	r.POST("/register", userHandler.Create)
	r.POST("/login", userHandler.Login)

	return r
}
