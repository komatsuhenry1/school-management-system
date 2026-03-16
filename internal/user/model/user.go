package model

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Cpf          string    `json:"cpf"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
