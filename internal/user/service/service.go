package service

import (
	"schoolmanagement/internal/user/model"
	"schoolmanagement/internal/user/repository"
)

type UserService interface {
	GetProfessionals() ([]model.User, error)
	UpdateUserPartial(id string, updates map[string]interface{}) (*model.User, error)
	// New methods for full CRUD
	GetAllUsers() ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
	DeleteUser(id string) error
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

// New methods for full CRUD

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepository.DeleteUser(id)
}
