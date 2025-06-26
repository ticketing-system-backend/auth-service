package service

import (
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/repository"
)

func GetAllRoles() ([]model.Role, error) {
	return repository.GetAllRole()
}

func GetRoleById(id uint) (*model.Role, error) {
	return repository.FindRoleById(id)
}

func CreateRole(role *model.Role) error {
	return repository.CreateRole(role)
}

func UpdateRole(role *model.Role) error {
	return repository.UpdateRole(role)
}
