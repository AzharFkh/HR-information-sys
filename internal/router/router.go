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

	employeeRoute := r.Use(middleware.AuthMiddleware())
	{
		employeeRoute.POST("/leaves", leaveHandler.CreateRequest)
		employeeRoute.GET("/leaves", leaveHandler.GetAllLeave)
	}

	hrRoute := r.Use(middleware.RequiredRole("hr"))
	
	{
		hrRoute.GET("/hr/pending", leaveHandler.GetPendingRequest)
		hrRoute.PATCH("/hr/update-status/:id", leaveHandler.UpdateStatus)
	}
	return r
}
