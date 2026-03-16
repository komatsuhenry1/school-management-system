package service

import (
	"schoolmanagement/internal/user/model"
	"schoolmanagement/internal/user/repository"
)

type UserService interface {
	GetProfessionals() ([]model.User, error)
	UpdateUserPartial(id string, updates map[string]interface{}) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetProfessionals() ([]model.User, error) {
	return s.userRepository.GetProfessionals()
}

func (s *userService) UpdateUserPartial(id string, updates map[string]interface{}) (*model.User, error) {
	return s.userRepository.UpdateUserPartial(id, updates)
}
