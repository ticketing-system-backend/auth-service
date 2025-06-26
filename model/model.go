package model

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	NamaLengkap string    `gorm:"not null" json:"nama_lengkap"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Password    string    `gorm:"not null" json:"-"`
	Roles       []Role    `gorm:"many2many:user_roles" json:"roles,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"not null" json:"nama"`
	Deskripsi string    `gorm:"not null" json:"deskripsi"`
	Level     string    `gorm:"type:role_level" json:"level"`
	Users     []User    `gorm:"many2many:user_roles" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}