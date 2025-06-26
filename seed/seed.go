package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/ticketing-system-backend/auth-service/config"
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/utils"
)

func SeedRoles() {
	roles := []model.Role{
		{Nama: "Admin", Deskripsi: "Admin sistem", Level: "admin"},
		{Nama: "Superadmin", Deskripsi: "Penguasa sistem", Level: "superadmin"},
		{Nama: "Staff", Deskripsi: "Staf operasional", Level: "staff"},
		{Nama: "HRD", Deskripsi: "Bagian SDM", Level: "hrd"},
		{Nama: "Finance", Deskripsi: "Bagian keuangan", Level: "finance"},
		{Nama: "Customer", Deskripsi: "Pengguna aplikasi", Level: "customer"},
	}

	for _, role := range roles {
		var existing model.Role
		err := config.DB.Where("level = ?", role.Level).First(&existing).Error
		if err != nil {
			role.CreatedAt = time.Now()
			role.UpdatedAt = time.Now()
			config.DB.Create(&role)
			fmt.Printf("Role '%s' ditambahkan\n", role.Level)
		}
	}
}

func SeedSuperAdmin() {
	var user model.User
	err := config.DB.Where("email = ?", "superadmin@example.com").Preload("Roles").First(&user).Error

	if err == nil {
		log.Println("Superadmin sudah ada")
		return
	}

	hashed, _ := utils.HashPassword("password")

	// Ambil role superadmin
	var superadminRole model.Role
	if err := config.DB.Where("level = ?", "superadmin").First(&superadminRole).Error; err != nil {
		fmt.Println("Role superadmin tidak ditemukan:", err)
		return
	}

	// Buat user superadmin dan attach role langsung
	user = model.User{
		NamaLengkap: "Superadmin",
		Email:       "superadmin@example.com",
		Password:    hashed,
		Roles:       []model.Role{superadminRole},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		fmt.Println("Gagal membuat user superadmin:", err)
		return
	}

	fmt.Println("Seeder superadmin berhasil dibuat")
}