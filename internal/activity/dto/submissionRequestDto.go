package dto

type ExerciseSubmissionDTO struct {
	ExerciseID    string `json:"exercise_id" binding:"required"`
	StudentAnswer string `json:"student_answer" binding:"required"`
}

type SubmissionRequestDTO struct {
	Answers    []ExerciseSubmissionDTO `json:"answers" binding:"required"`
}
