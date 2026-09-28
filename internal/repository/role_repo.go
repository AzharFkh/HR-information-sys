package repository

import (
	"hris-tools/internal/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByName(name string) (*models.Role, error){
	var role models.Role

	if err := r.db.First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	
	return &role, nil
}