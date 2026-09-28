package router

import (
	"hris-tools/internal/handler"
	"hris-tools/internal/middleware"
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
	
	r.POST("/register", userHandler.Create)
	r.POST("/login", userHandler.Login)

	leavesRoute := r.Use(middleware.AuthMiddleware())
	{
		leavesRoute.POST("/leaves", leaveHandler.CreateRequest)
		leavesRoute.GET("/leaves/pending", leaveHandler.GetPendingRequest)
	}

	return r
}
