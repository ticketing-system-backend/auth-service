package repository

import (
	"github.com/ticketing-system-backend/auth-service/config"
	"github.com/ticketing-system-backend/auth-service/model"
)

func FindRoleById(id uint) (*model.Role, error) {
	var role model.Role
	err := config.DB.First(&role, id).Error

	return &role, err
}

func GetAllRole() ([]model.Role, error) {
	var role []model.Role
	err := config.DB.Find(&role).Error

	return role, err
}

func CreateRole(role *model.Role) error {
	return config.DB.Create(role).Error
}

func UpdateRole(role *model.Role) error {
	return config.DB.Save(role).Error
}
