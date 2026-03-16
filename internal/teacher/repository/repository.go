package repository

import (
	"schoolmanagement/internal/user/model"

	"gorm.io/gorm"
)

type TeacherRepository interface {
	CreateTeacher(user *model.User) error
	GetAllTeachers() ([]model.User, error)
	GetTeacherByID(id string) (*model.User, error)
	UpdateTeacher(id string, updates map[string]interface{}) (*model.User, error)
	DeleteTeacher(id string) error
}

type teacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) CreateTeacher(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *teacherRepository) GetAllTeachers() ([]model.User, error) {
	var teachers []model.User
	if err := r.db.Where("role = ?", "TEACHER").Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *teacherRepository) GetTeacherByID(id string) (*model.User, error) {
	var teacher model.User
	if err := r.db.Where("id = ? AND role = ?", id, "TEACHER").First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) UpdateTeacher(id string, updates map[string]interface{}) (*model.User, error) {
	var teacher model.User
	if err := r.db.Where("id = ? AND role = ?", id, "TEACHER").First(&teacher).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&teacher).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) DeleteTeacher(id string) error {
	return r.db.Where("id = ? AND role = ?", id, "TEACHER").Delete(&model.User{}).Error
}
