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
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	leaveService := service.NewLeaveService(leaveRepo)
	userService := service.NewUserService(userRepo, roleRepo)

	leaveHandler := handler.NewLeaveHandler(leaveService)
	userHandler := handler.NewUserHandler(userService)

	r := router.Setuprouter(userHandler, leaveHandler)

	r.Run(":7100")
}
