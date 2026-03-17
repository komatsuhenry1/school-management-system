package dto

type HardestQuestionDTO struct {
	Question        string  `json:"question"`
	Subject         string  `json:"subject"`
	ErrorPercentage float32 `json:"error_percentage"`
}

type ActivityMetricsDTO struct {
	ClassAverage     float32 `json:"class_average"`
	HighestScore     float32 `json:"highest_score"`
	LowestScore      float32 `json:"lowest_score"`
	TotalSubmissions int     `json:"total_submissions"`
}

type StudentSubmissionStatusDTO struct {
	Name      string  `json:"name"`
	Submitted bool    `json:"submitted"`
	Score     float32 `json:"score"`
}

type ActivityDashboardDTO struct {
	Metrics          ActivityMetricsDTO           `json:"metrics"`
	HardestQuestions []HardestQuestionDTO         `json:"hardest_questions"`
	Students         []StudentSubmissionStatusDTO `json:"students"`
}

type SubjectAccuracyDTO struct {
	Subject  string  `json:"subject"`
	Accuracy float32 `json:"accuracy"`
}

type ActivityAccuracyDTO struct {
	ActivityID string  `json:"activity_id"`
	Title      string  `json:"title"`
	Accuracy   float32 `json:"accuracy"`
}

type StudentDashboardDTO struct {
	TotalActivitiesCompleted int                   `json:"total_activities_completed"`
	AverageScore             float32               `json:"average_score"`
	Subjects                 []SubjectAccuracyDTO  `json:"subjects"`
	Activities               []ActivityAccuracyDTO `json:"activities"`
}
