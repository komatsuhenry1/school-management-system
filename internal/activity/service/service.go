package service

import (
	"schoolmanagement/internal/activity/dto"
	"schoolmanagement/internal/activity/model"
	"schoolmanagement/internal/activity/repository"
	userRepo "schoolmanagement/internal/user/repository"
	"sort"
)

type ActivityService interface {
	CreateActivity(req *dto.ActivityRequestDTO) (*model.Activity, error)
	GetAllActivities() ([]model.Activity, error)
	GetActivityByID(id string) (*model.Activity, error)
	UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error)
	DeleteActivity(id string) error
	SubmitActivity(req *dto.SubmissionRequestDTO) (*model.ActivitySubmission, error)
	GetActivityDashboard(activityID string) (*dto.ActivityDashboardDTO, error)
}

type activityService struct {
	activityRepository repository.ActivityRepository
	userRepository     userRepo.UserRepository
}

func NewActivityService(activityRepository repository.ActivityRepository, userRepository userRepo.UserRepository) ActivityService {
	return &activityService{
		activityRepository: activityRepository,
		userRepository:     userRepository,
	}
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

func (s *activityService) GetActivityDashboard(activityID string) (*dto.ActivityDashboardDTO, error) {
	// 1. Fetch the activity itself to know exercises
	activity, err := s.activityRepository.GetActivityByID(activityID)
	if err != nil {
		return nil, err
	}

	exercisesMap := make(map[string]model.Exercise)
	for _, ex := range activity.Exercises {
		exercisesMap[ex.ID] = ex
	}

	// 2. Fetch all submissions for this activity
	submissions, err := s.activityRepository.GetSubmissionsByActivityID(activityID)
	if err != nil {
		return nil, err
	}

	// 3. Fetch all students (role = USER)
	students, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Metrics
	var totalScore float32 = 0
	var highestScore float32 = 0
	lowestScore := float32(1000) // Arbitrary high number
	if len(submissions) == 0 {
		lowestScore = 0
	}

	// Error tracking per exercise: map[ExerciseID]int (count of errors)
	exerciseErrors := make(map[string]int)
	totalAttemptsPerExercise := make(map[string]int)

	// Map submissions by UserID for quick student status check
	submissionsByUser := make(map[string]*model.ActivitySubmission)

	for i := range submissions {
		sub := &submissions[i]
		submissionsByUser[sub.UserID] = sub

		// Aggregate scores
		totalScore += sub.Score
		if sub.Score > highestScore {
			highestScore = sub.Score
		}
		if sub.Score < lowestScore {
			lowestScore = sub.Score
		}

		// Aggregate exercise errors
		for _, ans := range sub.Answers {
			totalAttemptsPerExercise[ans.ExerciseID]++
			if !ans.IsCorrect {
				exerciseErrors[ans.ExerciseID]++
			}
		}
	}

	var classAverage float32 = 0
	if len(submissions) > 0 {
		classAverage = totalScore / float32(len(submissions))
	}

	// 4. Calculate Hardest Questions
	var hardestQuestions []dto.HardestQuestionDTO
	for exID, errCount := range exerciseErrors {
		attempts := totalAttemptsPerExercise[exID]
		if attempts > 0 {
			errPct := (float32(errCount) / float32(attempts)) * 100
			ex, found := exercisesMap[exID]
			if found {
				hardestQuestions = append(hardestQuestions, dto.HardestQuestionDTO{
					Question:        ex.Question,
					Subject:         ex.ExerciseSubject,
					ErrorPercentage: errPct,
				})
			}
		}
	}

	// Sort hardest questions descending by ErrorPercentage
	sort.Slice(hardestQuestions, func(i, j int) bool {
		return hardestQuestions[i].ErrorPercentage > hardestQuestions[j].ErrorPercentage
	})

	// Keep only top 3
	if len(hardestQuestions) > 3 {
		hardestQuestions = hardestQuestions[:3]
	}

	// 5. Build Student List
	var studentStatuses []dto.StudentSubmissionStatusDTO
	for _, student := range students {
		status := dto.StudentSubmissionStatusDTO{
			Name:      student.Name,
			Submitted: false,
			Score:     0,
		}

		if sub, found := submissionsByUser[student.ID]; found {
			status.Submitted = true
			status.Score = sub.Score
		}

		studentStatuses = append(studentStatuses, status)
	}

	// 6. Build final DTO
	dashboardDTO := &dto.ActivityDashboardDTO{
		Metrics: dto.ActivityMetricsDTO{
			ClassAverage:     classAverage,
			HighestScore:     highestScore,
			LowestScore:      lowestScore,
			TotalSubmissions: len(submissions),
		},
		HardestQuestions: hardestQuestions,
		Students:         studentStatuses,
	}

	return dashboardDTO, nil
}
