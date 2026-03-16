package model

type Alternative struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ExerciseID string `gorm:"type:uuid;not null;index" json:"exercise_id"` // Foreign Key
	Letter     string `json:"letter"`
	Value      string `json:"value"`
}
