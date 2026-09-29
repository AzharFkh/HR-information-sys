package main

import (
	"hris-tools/cmd/seed"
	"hris-tools/internal/config"
	"hris-tools/internal/models"
	"log"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.LeaveType{},
		&models.LeaveRequest{},
		&models.LeaveBalance{},
	)

	if err != nil {
		log.Fatal("migratioin failed: ", err)
	}

	seed.SetupAndSeedDatabase(config.DB)

	log.Println("migration success and seeding")
}
