package service

import (
	"errors"
	"schoolmanagement/internal/activity/dto"
	"schoolmanagement/internal/activity/model"
	"schoolmanagement/internal/activity/repository"
	userRepo "schoolmanagement/internal/user/repository"
	"sort"
	"fmt"
)

type ActivityService interface {
	CreateActivity(req *dto.ActivityRequestDTO) (*model.Activity, error)
	GetAllActivities() ([]model.Activity, error)
	GetActivityByID(id string) (*model.Activity, error)
	UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error)
	DeleteActivity(id string) error
	SubmitActivity(req *dto.SubmissionRequestDTO, activityID string, userID string) (*model.ActivitySubmission, error)
	GetActivityDashboard(activityID string) (*dto.ActivityDashboardDTO, error)
	GetActiveActivities(userID string) ([]dto.ActiveActivityResponseDTO, error)
	GetActivityQuestions(activityID string) (*dto.ActivityQuestionsResponseDTO, error)
	GetStudentDashboard(userID string) (*dto.StudentDashboardDTO, error)
	GetClassRanking() ([]dto.StudentRankingDTO, error)
	GetClassroomMetrics() (*dto.ClassroomMetricsDTO, error)
	UpdateAlternative(alternativeID string, updates map[string]interface{}) error
	UpdateActivityFull(id string, req *dto.ActivityRequestDTO) (*model.Activity, error)
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
		Subject:       req.Subject,
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

func (s *activityService) SubmitActivity(req *dto.SubmissionRequestDTO, activityID string, userID string) (*model.ActivitySubmission, error) {
	// 0. Check for duplicate submission
	hasSubmitted, err := s.activityRepository.HasUserSubmittedActivity(userID, activityID)
	if err != nil {
		return nil, err
	}
	if hasSubmitted {
		return nil, errors.New("o usuário já submeteu uma resposta para esta atividade")
	}

	// 1. Fetch the activity to know the correct answers
	activity, err := s.activityRepository.GetActivityByID(activityID)
	if err != nil {
		return nil, err
	}

	// 2. Map correct answers for O(1) lookup-
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
		UserID:     userID,
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

	fmt.Println("========")
	fmt.Println("========")
	fmt.Println(hardestQuestions)
	fmt.Println("========")

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

func (s *activityService) GetActiveActivities(userID string) ([]dto.ActiveActivityResponseDTO, error) {
	activities, err := s.activityRepository.GetActiveActivities()
	if err != nil {
		return nil, err
	}

	// Fetch user's submissions to check which ones are submitted
	userSubmissions := make(map[string]float32)
	if userID != "" {
		submissions, err := s.activityRepository.GetSubmissionsByUserID(userID)
		if err == nil {
			for _, sub := range submissions {
				userSubmissions[sub.ActivityID] = sub.Score
			}
		}
	}

	var response []dto.ActiveActivityResponseDTO
	for _, act := range activities {
		var exercises []dto.StudentExerciseDTO
		for _, ex := range act.Exercises {
			// Map alternatives if they exist
			var alts []dto.AlternativeDTO
			for _, a := range ex.Alternatives {
				alts = append(alts, dto.AlternativeDTO{
					Letter: a.Letter,
					Value:  a.Value,
				})
			}

			exercises = append(exercises, dto.StudentExerciseDTO{
				ID:              ex.ID,
				ExerciseNumber:  ex.ExerciseNumber,
				ExerciseSubject: ex.ExerciseSubject,
				Question:        ex.Question,
				ExerciseValue:   ex.ExerciseValue,
				Alternatives:    alts,
			})
		}

		isSubmitted := false
		var scorePtr *float32

		if score, exists := userSubmissions[act.ID]; exists {
			isSubmitted = true
			scoreCopy := score
			scorePtr = &scoreCopy
		}

		response = append(response, dto.ActiveActivityResponseDTO{
			ID:            act.ID,
			Title:         act.Title,
			Description:   act.Description,
			ActivityValue: act.ActivityValue,
			Subject:       act.Subject,
			Status:        act.Status,
			IsSubmitted:   isSubmitted,
			Score:         scorePtr,
			Exercises:     exercises,
			CreatedAt:     act.CreatedAt,
		})
	}

	return response, nil
}

func (s *activityService) GetActivityQuestions(activityID string) (*dto.ActivityQuestionsResponseDTO, error) {
	activity, err := s.activityRepository.GetActivityByID(activityID)
	if err != nil {
		return nil, err
	}

	var exercises []dto.StudentExerciseDTO
	for _, ex := range activity.Exercises {
		// Map alternatives
		var alts []dto.AlternativeDTO
		for _, a := range ex.Alternatives {
			alts = append(alts, dto.AlternativeDTO{
				Letter: a.Letter,
				Value:  a.Value,
			})
		}

		exercises = append(exercises, dto.StudentExerciseDTO{
			ID:              ex.ID,
			ExerciseNumber:  ex.ExerciseNumber,
			ExerciseSubject: ex.ExerciseSubject,
			Question:        ex.Question,
			ExerciseValue:   ex.ExerciseValue,
			Alternatives:    alts, // Include alternatives but NO Answer
		})
	}

	response := &dto.ActivityQuestionsResponseDTO{
		ID:            activity.ID,
		Title:         activity.Title,
		Description:   activity.Description,
		ActivityValue: activity.ActivityValue,
		Exercises:     exercises,
	}

	return response, nil
}

func (s *activityService) GetStudentDashboard(userID string) (*dto.StudentDashboardDTO, error) {
	// 1. Fetch student's submissions
	submissions, err := s.activityRepository.GetSubmissionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	totalActivitiesCompleted := len(submissions)
	if totalActivitiesCompleted == 0 {
		return &dto.StudentDashboardDTO{
			TotalActivitiesCompleted: 0,
			AverageScore:             0,
			Subjects:                 []dto.SubjectAccuracyDTO{},
			Activities:               []dto.ActivityAccuracyDTO{},
		}, nil
	}

	var totalScore float32 = 0
	
	// Track correct anwers vs total answers per subject
	subjectTotalAnswers := make(map[string]int)
	subjectCorrectAnswers := make(map[string]int)

	// Track correct answers vs total answers per activity subject
	activitySubTotalAnswers := make(map[string]int)
	activitySubCorrectAnswers := make(map[string]int)

	// Since we need the subject, we have to fetch the exercises related to the answers.
	// Instead of querying exercise by exercise, we will collect all unique Activity IDs
	// from the submissions, fetch those activities, and build an Exercise map.
	activityIDs := make(map[string]bool)
	for _, sub := range submissions {
		totalScore += sub.Score
		activityIDs[sub.ActivityID] = true
	}

	averageScore := totalScore / float32(totalActivitiesCompleted)

	// Fetch all necessary exercises to get subjects and map activities for Titles
	exercisesMap := make(map[string]model.Exercise)
	activitiesMap := make(map[string]model.Activity)
	for actID := range activityIDs {
		act, err := s.activityRepository.GetActivityByID(actID)
		if err == nil {
			activitiesMap[act.ID] = *act
			for _, ex := range act.Exercises {
				exercisesMap[ex.ID] = ex
			}
		}
	}

	// We will use the Subject field directly from the Activity for `activities` accuracy
	activitySubjects := make(map[string]string)
	for actID, act := range activitiesMap {
		if act.Subject != "" {
			activitySubjects[actID] = act.Subject
		} else {
			activitySubjects[actID] = "Geral"
		}
	}

	for _, sub := range submissions {
		actSubject := activitySubjects[sub.ActivityID]

		for _, ans := range sub.Answers {
			// 1. For activities array (matéria da atividade)
			activitySubTotalAnswers[actSubject]++
			if ans.IsCorrect {
				activitySubCorrectAnswers[actSubject]++
			}

			// 2. For subjects array (matéria do exercício individual)
			if ex, exists := exercisesMap[ans.ExerciseID]; exists {
				subjectTotalAnswers[ex.ExerciseSubject]++
				if ans.IsCorrect {
					subjectCorrectAnswers[ex.ExerciseSubject]++
				}
			}
		}
	}

	var subjectsAccuracy []dto.SubjectAccuracyDTO
	for subject, totalAnswers := range subjectTotalAnswers {
		correctAnswers := subjectCorrectAnswers[subject]
		accuracy := (float32(correctAnswers) / float32(totalAnswers)) * 100

		subjectsAccuracy = append(subjectsAccuracy, dto.SubjectAccuracyDTO{
			Subject:  subject,
			Accuracy: accuracy,
		})
	}

	var activitiesAccuracy []dto.ActivityAccuracyDTO
	for actSubj, totalAnswers := range activitySubTotalAnswers {
		correctAnswers := activitySubCorrectAnswers[actSubj]
		accuracy := (float32(correctAnswers) / float32(totalAnswers)) * 100

		activitiesAccuracy = append(activitiesAccuracy, dto.ActivityAccuracyDTO{
			Subject:    actSubj,
			Accuracy:   accuracy,
		})
	}

	// Sort subjects alphabetically or by accuracy if preferred, leaving as is for now

	dashboard := &dto.StudentDashboardDTO{
		TotalActivitiesCompleted: totalActivitiesCompleted,
		AverageScore:             averageScore,
		Subjects:                 subjectsAccuracy,
		Activities:               activitiesAccuracy,
	}

	return dashboard, nil
}

func (s *activityService) GetClassRanking() ([]dto.StudentRankingDTO, error) {
	// 1. Fetch all students
	students, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// 2. Fetch all submissions
	submissions, err := s.activityRepository.GetAllSubmissions()
	if err != nil {
		return nil, err
	}

	// 3. Aggregate scores per UserID
	userScores := make(map[string]float32)
	for _, sub := range submissions {
		userScores[sub.UserID] += sub.Score
	}

	// 4. Map user scores to ranking DTOs using student names
	var rankings []dto.StudentRankingDTO
	for _, student := range students {
		// Even if they have no submissions, we can include them with 0 score, 
		// but since we want the top ranking, we'll sort them down.
		rankings = append(rankings, dto.StudentRankingDTO{
			Name:  student.Name,
			Score: userScores[student.ID],
		})
	}

	// 5. Sort rankings descending by score
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	// 6. Assign positions and filter top 3
	var topRankings []dto.StudentRankingDTO
	limit := 3
	if len(rankings) < limit {
		limit = len(rankings)
	}

	for i := 0; i < limit; i++ {
		// Only include students who have at least > 0 score? 
		// Or all top 3 even if score is 0? Let's just include the top 3.
		if rankings[i].Score > 0 {
			rankings[i].Position = i + 1
			topRankings = append(topRankings, rankings[i])
		}
	}

	// If no one has > 0 score, returns empty slice instead of nil for better JSON
	if topRankings == nil {
		topRankings = []dto.StudentRankingDTO{}
	}

	return topRankings, nil
}

func (s *activityService) GetClassroomMetrics() (*dto.ClassroomMetricsDTO, error) {
	// 1. Fetch data
	students, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	activities, err := s.activityRepository.GetAllActivities()
	if err != nil {
		return nil, err
	}

	submissions, err := s.activityRepository.GetAllSubmissions()
	if err != nil {
		return nil, err
	}

	// 2. Build exercises map for grouping subjects
	exercisesMap := make(map[string]model.Exercise)
	for _, act := range activities {
		for _, ex := range act.Exercises {
			exercisesMap[ex.ID] = ex
		}
	}

	// 3. Aggregate metrics
	var totalScore float32 = 0
	var totalCorrectAnswers int = 0
	var totalAnswersCount int = 0

	subjectTotalAnswers := make(map[string]int)
	subjectCorrectAnswers := make(map[string]int)

	for _, sub := range submissions {
		totalScore += sub.Score

		for _, ans := range sub.Answers {
			totalAnswersCount++
			if ans.IsCorrect {
				totalCorrectAnswers++
			}

			if ex, exists := exercisesMap[ans.ExerciseID]; exists {
				subjectTotalAnswers[ex.ExerciseSubject]++
				if ans.IsCorrect {
					subjectCorrectAnswers[ex.ExerciseSubject]++
				}
			}
		}
	}

	var classAverageScore float32 = 0
	if len(submissions) > 0 {
		classAverageScore = totalScore / float32(len(submissions))
	}

	var generalAccuracy float32 = 0
	if totalAnswersCount > 0 {
		generalAccuracy = (float32(totalCorrectAnswers) / float32(totalAnswersCount)) * 100
	}

	// 4. Calculate Top 3 Hardest Subjects (lowest accuracy)
	var subjectAccuracies []dto.SubjectAccuracyDTO
	for subj, totalAns := range subjectTotalAnswers {
		if totalAns > 0 {
			acc := (float32(subjectCorrectAnswers[subj]) / float32(totalAns)) * 100
			subjectAccuracies = append(subjectAccuracies, dto.SubjectAccuracyDTO{
				Subject:  subj,
				Accuracy: acc,
			})
		}
	}

	sort.Slice(subjectAccuracies, func(i, j int) bool {
		return subjectAccuracies[i].Accuracy < subjectAccuracies[j].Accuracy
	})

	var hardestSubjects []dto.SubjectAccuracyDTO
	limit := 3
	if len(subjectAccuracies) < limit {
		limit = len(subjectAccuracies)
	}

	for i := 0; i < limit; i++ {
		hardestSubjects = append(hardestSubjects, subjectAccuracies[i])
	}

	if hardestSubjects == nil {
		hardestSubjects = []dto.SubjectAccuracyDTO{}
	}

	return &dto.ClassroomMetricsDTO{
		TotalStudents:     len(students),
		TotalActivities:   len(activities),
		TotalSubmissions:  len(submissions),
		ClassAverageScore: classAverageScore,
		GeneralAccuracy:   generalAccuracy,
		HardestSubjects:   hardestSubjects,
	}, nil
}

func (s *activityService) UpdateAlternative(alternativeID string, updates map[string]interface{}) error {
	// For better security, you could also fetch the activity and check if the alternative 
	// actually belongs to the provided activity ID in the route, but in a simple implementation
	// we just update by the unique alternative ID.
	return s.activityRepository.UpdateAlternative(alternativeID, updates)
}

func (s *activityService) UpdateActivityFull(id string, req *dto.ActivityRequestDTO) (*model.Activity, error) {
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
			ActivityID:      id,
			ExerciseNumber:  e.ExerciseNumber,
			ExerciseSubject: e.ExerciseSubject,
			Question:        e.Question,
			Answer:          e.Answer,
			ExerciseValue:   e.ExerciseValue,
			Alternatives:    alts,
		})
	}

	activity := &model.Activity{
		ID:            id,
		Title:         req.Title,
		Description:   req.Description,
		Subject:       req.Subject,
		ActivityValue: req.ActivityValue,
		Exercises:     exercises,
	}

	if err := s.activityRepository.UpdateActivityFull(activity); err != nil {
		return nil, err
	}
	return activity, nil
}
