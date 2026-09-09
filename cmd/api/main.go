package main

import (
	"hris-tools/internal/config"
	"hris-tools/internal/handler"
	"hris-tools/internal/repository"
	"hris-tools/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main(){

	config.LoadEnv()
	config.ConnectDB()

	db := config.DB

	leaveRepo := repository.NewLeavesRepository(db)

	leaveService := service.NewLeaveService(leaveRepo)

	leaveHandler := handler.NewLeaveHandler(leaveService)


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "uji",
		})
	})

	r.POST("/leaves", leaveHandler.CreateRequest)

	r.GET("/leaves/pending", leaveHandler.GetPendingRequest)

	r.Run(":7100")
}