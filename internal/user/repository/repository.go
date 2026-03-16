package repository

import (
	"schoolmanagement/internal/user/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUserNameOrEmail(name, email string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	GetProfessionals() ([]model.User, error)
	UpdateUserPartial(id string, updates map[string]interface{}) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUserNameOrEmail(name, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("name = ? OR email = ?", name, email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) GetProfessionals() ([]model.User, error) {
	var professionals []model.User
	if err := r.db.Where("role = ?", "PROFESSIONAL").Find(&professionals).Error; err != nil {
		return nil, err
	}
	return professionals, nil
}

func (r *userRepository) UpdateUserPartial(id string, updates map[string]interface{}) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
