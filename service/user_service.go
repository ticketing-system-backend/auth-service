package service

import (
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/repository"
	"github.com/ticketing-system-backend/auth-service/utils"
)

func GetAllUsers() ([]model.User, error) {
	return repository.GetAllUser()
}

func GetUserById(id uint) (*model.User, error) {
	return repository.FindUserById(id)
}

func CreateUser(user *model.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return repository.CreateUser(user)
}

func UpdateUser(user *model.User) error {
	if user.Password != "" {
		hashed, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}
	return repository.UpdateUser(user)
}
