package service

import (
	"schoolmanagement/internal/teacher/dto"
	"schoolmanagement/internal/teacher/repository"
	"schoolmanagement/internal/user/model"
)

type TeacherService interface {
	CreateTeacher(req *dto.TeacherRequestDTO) (*model.User, error)
	GetAllTeachers() ([]model.User, error)
	GetTeacherByID(id string) (*model.User, error)
	UpdateTeacher(id string, updates map[string]interface{}) (*model.User, error)
	DeleteTeacher(id string) error
}

type teacherService struct {
	teacherRepository repository.TeacherRepository
}

func NewTeacherService(teacherRepository repository.TeacherRepository) TeacherService {
	return &teacherService{teacherRepository: teacherRepository}
}

func (s *teacherService) CreateTeacher(req *dto.TeacherRequestDTO) (*model.User, error) {
	teacher := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Role:  "TEACHER",
	}
	if err := s.teacherRepository.CreateTeacher(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *teacherService) GetAllTeachers() ([]model.User, error) {
	return s.teacherRepository.GetAllTeachers()
}

func (s *teacherService) GetTeacherByID(id string) (*model.User, error) {
	return s.teacherRepository.GetTeacherByID(id)
}

func (s *teacherService) UpdateTeacher(id string, updates map[string]interface{}) (*model.User, error) {
	return s.teacherRepository.UpdateTeacher(id, updates)
}

func (s *teacherService) DeleteTeacher(id string) error {
	return s.teacherRepository.DeleteTeacher(id)
}
