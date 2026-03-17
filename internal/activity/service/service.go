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
	SubmitActivity(req *dto.SubmissionRequestDTO) (*model.ActivitySubmission, error)
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

func (s *activityService) SubmitActivity(req *dto.SubmissionRequestDTO) (*model.ActivitySubmission, error) {
	// 1. Fetch the activity to know the correct answers
	activity, err := s.activityRepository.GetActivityByID(req.ActivityID)
	if err != nil {
		return nil, err
	}

	// 2. Map correct answers for O(1) lookup
	correctAnswers := make(map[string]model.Exercise)
	for _, ex := range activity.Exercises {
		correctAnswers[ex.ID] = ex
	}

	var totalScore float32
	var exerciseSubmissions []model.ExerciseSubmission

	// 3. Process each student answer
	for _, answerReq := range req.Answers {
		exercise, exists := correctAnswers[answerReq.ExerciseID]
		if !exists {
			// Skip or return error if they answered an exercise that doesn't belong to this activity
			continue
		}

		isCorrect := false
		pointsEarned := float32(0)

		// Simple exact match logic (could be improved for case-insensitivity or alternatives in the future)
		if answerReq.StudentAnswer == exercise.Answer {
			isCorrect = true
			pointsEarned = exercise.ExerciseValue
			totalScore += pointsEarned
		}

		exerciseSubmissions = append(exerciseSubmissions, model.ExerciseSubmission{
			ExerciseID:    answerReq.ExerciseID,
			StudentAnswer: answerReq.StudentAnswer,
			IsCorrect:     isCorrect,
			PointsEarned:  pointsEarned,
		})
	}

	// 4. Create the final submission payload
	submission := &model.ActivitySubmission{
		ActivityID: activity.ID,
		UserID:     req.UserID,
		Score:      totalScore,
		Status:     "COMPLETED", // Adjust as necessary if manual review is needed
		Answers:    exerciseSubmissions,
	}

	// 5. Save to database
	if err := s.activityRepository.SubmitActivity(submission); err != nil {
		return nil, err
	}

	return submission, nil
}
