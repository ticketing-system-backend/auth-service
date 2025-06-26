package service

import (
	"errors"

	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/repository"
	"github.com/ticketing-system-backend/auth-service/utils"
)

func Login(email, password string, isDashboard bool) (*model.User, error) {
	user, err := repository.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("email tidak ditemukan")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("password salah")
	}

	if isDashboard {
		// tidak boleh login jika hanya punya role customer
		isAllowed := false
		for _, ur := range user.Roles {
			if ur.Level != "customer" {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return nil, errors.New("akses dashboard ditolak")
		}
	}

	return user, nil
}
