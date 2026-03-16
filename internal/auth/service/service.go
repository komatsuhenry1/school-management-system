package service

import (
	"schoolmanagement/internal/user/repository"
	"schoolmanagement/internal/auth/dto"
	"schoolmanagement/internal/utils"
	"schoolmanagement/internal/user/model"
	"errors"
	"fmt"
	"strings"
)

type UserService interface {
	RegisterUser(userRequestDTO *dto.UserRequestDTO) error
	LoginUser(requestDto *dto.LoginRequestDTO) (string, string, model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) RegisterUser(userRequestDTO *dto.UserRequestDTO) error {
	normalizedEmail, err := utils.EmailRegex(userRequestDTO.Email)
	if err != nil {
		return err
	}

	normalizedCPF, err := utils.ValidateCPF(userRequestDTO.Cpf)
	if err != nil {
		return err
	}

	normalizedPhone, err := utils.ValidatePhone(userRequestDTO.Phone)
	if err != nil {
		return err
	}

	// normalizedCEP, err := utils.ValidateCEP(userRequestDTO.CEP)
	// if err != nil {
	// 	return err
	// }

	if err := utils.ValidatePasswordRegex(userRequestDTO.Password); err != nil {
		return err
	}

	userRequestDTO.Name = strings.ToLower(userRequestDTO.Name)
	existingUser, err := s.userRepository.GetUserByUserNameOrEmail(userRequestDTO.Name, userRequestDTO.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("usuário já existe")
	}

	fmt.Println(userRequestDTO.Password)
	err = utils.HashPassword(&userRequestDTO.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:         utils.CapitalizeWords(userRequestDTO.Name),
		Email:        normalizedEmail,
		Cpf:          normalizedCPF,
		Password:     userRequestDTO.Password,
		Phone:        normalizedPhone,
		Role:         "USER",
	}

	err = s.userRepository.CreateUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) LoginUser(loginRequestDto *dto.LoginRequestDTO) (string, string, model.User, error) {
	loginRequestDto.Email = strings.ToLower(loginRequestDto.Email)
	user, err := s.userRepository.GetUserByEmail(loginRequestDto.Email)
	if err != nil {
		return "", "", model.User{}, fmt.Errorf("usuário ou senha incorretos")
	}

	if !utils.ComparePassword(user.Password, loginRequestDto.Password) {
		return "", "", model.User{}, fmt.Errorf("usuário ou senha incorretos")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", model.User{}, fmt.Errorf("erro ao gerar token")
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", model.User{}, fmt.Errorf("erro ao gerar refresh token")
	}

	// salvar o refresh token no banco de dados
	user.RefreshToken = utils.HashToken(refreshToken)
	err = s.userRepository.UpdateUser(user)
	if err != nil {
		return "", "", model.User{}, fmt.Errorf("erro ao salvar refresh token")
	}

	return token, refreshToken, *user, nil
}
