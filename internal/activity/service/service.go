package service

import (
	"schoolmanagement/internal/activity/dto"
	"schoolmanagement/internal/activity/model"
	"schoolmanagement/internal/activity/repository"
)

type ActivityService interface {
	CreateActivity(req *dto.ActivityRequestDTO) (*model.Activity, error)
	GetAllActivities() ([]model.Activity, error)
	GetActivityByID(id string) (*model.Activity, error)
	UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error)
	DeleteActivity(id string) error
}

type activityService struct {
	activityRepository repository.ActivityRepository
}

func NewActivityService(activityRepository repository.ActivityRepository) ActivityService {
	return &activityService{activityRepository: activityRepository}
}

func (s *activityService) CreateActivity(req *dto.ActivityRequestDTO) (*model.Activity, error) {
	exercises := make([]model.Exercise, 0, len(req.Exercises))
	
	for _, e := range req.Exercises {
		alts := make([]model.Alternative, 0, len(e.Alternatives))
		for _, alt := range e.Alternatives {
			alts = append(alts, model.Alternative{
				Letter: alt.Letter,
				Value:  alt.Value,
			})
		}

		exercises = append(exercises, model.Exercise{
			ExerciseNumber:  e.ExerciseNumber,
			ExerciseSubject: e.ExerciseSubject,
			Question:        e.Question,
			Answer:          e.Answer,
			ExerciseValue:   e.ExerciseValue,
			Alternatives:    alts,
		})
	}

	activity := &model.Activity{
		Title:         req.Title,
		Description:   req.Description,
		ActivityValue: req.ActivityValue,
		Exercises:     exercises,
	}

	if err := s.activityRepository.CreateActivity(activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *activityService) GetAllActivities() ([]model.Activity, error) {
	return s.activityRepository.GetAllActivities()
}

func (s *activityService) GetActivityByID(id string) (*model.Activity, error) {
	return s.activityRepository.GetActivityByID(id)
}

func (s *activityService) UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error) {
	return s.activityRepository.UpdateActivity(id, updates)
}

func (s *activityService) DeleteActivity(id string) error {
	return s.activityRepository.DeleteActivity(id)
}
