package repository

import (
	"github.com/ticketing-system-backend/auth-service/config"
	"github.com/ticketing-system-backend/auth-service/model"
)

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := config.DB.Preload("Roles").Where("email = ?", email).First(&user).Error

	return &user, err
}

func FindUserById(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.Preload("Roles").First(&user, id).Error

	return &user, err
}

func GetAllUser() ([]model.User, error) {
	var user []model.User
	err := config.DB.Preload("Roles").Find(&user).Error

	return user, err
}

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func UpdateUser(user *model.User) error {
	return config.DB.Save(user).Error
}
