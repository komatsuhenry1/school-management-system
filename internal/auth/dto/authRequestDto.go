package dto

import "errors"

type UserRequestDTO struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Cpf      string `json:"cpf" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

func (u *UserRequestDTO) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.Phone == "" {
		return errors.New("phone is required")
	}
	return nil
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
