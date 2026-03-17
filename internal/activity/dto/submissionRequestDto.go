package dto

type ExerciseSubmissionDTO struct {
	ExerciseID    string `json:"exercise_id" binding:"required"`
	StudentAnswer string `json:"student_answer" binding:"required"`
}

type SubmissionRequestDTO struct {
	ActivityID string                  `json:"activity_id" binding:"required"`
	UserID     string                  `json:"user_id" binding:"required"`
	Answers    []ExerciseSubmissionDTO `json:"answers" binding:"required"`
}
