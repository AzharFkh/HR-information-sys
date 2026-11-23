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
	// Public routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "uji",
		})
	})

	r.POST("/register", userHandler.Create)
	r.POST("/login", userHandler.Login)

	// Protected routes (authenticated users)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/leaves", leaveHandler.CreateRequest)
		protected.GET("/leaves", leaveHandler.GetAllLeave)

		// HR-only routes
		hr := protected.Group("/hr")
		hr.Use(middleware.RequiredRole("hr"))
		{
			hr.GET("/pending", leaveHandler.GetPendingRequest)
			hr.PATCH("/update-status/:id", leaveHandler.UpdateStatus)
		}
	}

	return r
}
