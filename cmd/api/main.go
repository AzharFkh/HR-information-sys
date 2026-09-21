package main

import (
	"hris-tools/internal/config"
	"hris-tools/internal/handler"
	"hris-tools/internal/repository"
	"hris-tools/internal/router"
	"hris-tools/internal/service"
)

func main() {

	config.LoadEnv()
	config.ConnectDB()

	db := config.DB

	leaveRepo := repository.NewLeavesRepository(db)

	leaveService := service.NewLeaveService(leaveRepo)

	leaveHandler := handler.NewLeaveHandler(leaveService)

	r := router.Setuprouter(leaveHandler)

	r.Run(":7100")
}
