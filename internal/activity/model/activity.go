package model

import (
	"time"
)

type Activity struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ActivityValue float32    `json:"activity_value"`
	Status        string     `json:"status"` // ACTIVE, INACTIVE
	Exercises     []Exercise `gorm:"foreignKey:ActivityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercises"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
