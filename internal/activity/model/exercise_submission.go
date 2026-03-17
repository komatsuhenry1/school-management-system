package model

type ExerciseSubmission struct {
	ID            string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SubmissionID      string  `gorm:"type:uuid;not null;index" json:"submission_id"` // Foreign Key
	ExerciseID    string  `gorm:"type:uuid;not null;index" json:"exercise_id"`   // Reference to the original exercise
	StudentAnswer string  `json:"student_answer"`
	IsCorrect     bool    `json:"is_correct"`
	PointsEarned  float32 `json:"points_earned"`
}
