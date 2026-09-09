package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Password string // harus di hashed
	// bisa tambahkan employee ID 

	RoleID uint `json:"role_id"`
	Role   Role `gorm:"foreignKey:RoleID"`
}
