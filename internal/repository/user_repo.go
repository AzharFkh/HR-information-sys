package repository

import (
	"errors"
	"hris-tools/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	UserByID(id uint) (*models.User, error)
	UserByName(name string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// create new user
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// find user by ID
func (r *userRepository) UserByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &user, nil
}

// find user by name
func (r *userRepository) UserByName(name string) (*models.User, error) {
	var user models.User

	if err := r.db.Preload("Role").First(&user, "name = ?", name).Error 
		err != nil{
			return nil, err
		}

	return &user, nil
}
