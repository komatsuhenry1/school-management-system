package model

import (
	"time"
)

type ActivitySubmission struct {
	ID        string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ActivityID  string               `gorm:"type:uuid;not null;index" json:"activity_id"`
	UserID    string               `gorm:"type:uuid;not null;index" json:"user_id"`
	Score     float32              `json:"score"`
	Status    string               `json:"status"` // e.g., COMPLETED, PENDING_REVIEW
	Answers   []ExerciseSubmission `gorm:"foreignKey:SubmissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"answers"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
