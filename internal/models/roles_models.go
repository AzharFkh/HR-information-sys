package models

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"` // admin, user, guest
}
