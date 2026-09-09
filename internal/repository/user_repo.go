package repository

import (
	"errors"
	"hris-tools/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	UserByID(id uint) (*models.User, error)
	// UserByName(name string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UserByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &user, nil
}
