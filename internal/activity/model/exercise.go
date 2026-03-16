package model

type Exercise struct {
	ID              string        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ActivityID      string        `gorm:"type:uuid;not null;index" json:"activity_id"` // Foreign Key
	ExerciseNumber  int           `json:"exercise_number"`
	ExerciseSubject string        `json:"exercise_subject"`
	Question        string        `json:"question"`
	Answer          string        `json:"answer"`
	ExerciseValue   float32       `json:"exercise_value"`
	Alternatives    []Alternative `gorm:"foreignKey:ExerciseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"alternatives"`
}
