package dto

import "time"

type AlternativeDTO struct {
	Letter string `json:"letter"`
	Value  string `json:"value"`
}

type ExerciseDTO struct {
	ExerciseNumber  int              `json:"exercise_number"`
	ExerciseSubject string           `json:"exercise_subject"`
	Question        string           `json:"question"`
	Answer          string           `json:"answer"`
	ExerciseValue   float32          `json:"exercise_value"`
	Alternatives    []AlternativeDTO `json:"alternatives"`
}

type ActivityRequestDTO struct {
	Title         string        `json:"title" binding:"required"`
	Description   string        `json:"description"`
	Subject       string        `json:"subject" binding:"required"`
	ActivityValue float32       `json:"activity_value"`
	Exercises     []ExerciseDTO `json:"exercises"`
}

type StudentExerciseDTO struct {
	ID              string           `json:"id"`
	ExerciseNumber  int              `json:"exercise_number"`
	ExerciseSubject string           `json:"exercise_subject"`
	Question        string           `json:"question"`
	ExerciseValue   float32          `json:"exercise_value"`
	Alternatives    []AlternativeDTO `json:"alternatives"`
}

type ActiveActivityResponseDTO struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	ActivityValue float32              `json:"activity_value"`
	Subject       string               `json:"subject"`
	Status        string               `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	IsSubmitted   bool                 `json:"is_submitted"`
	Score         *float32             `json:"score,omitempty"`
	Exercises     []StudentExerciseDTO `json:"exercises"`
}

type ActivityQuestionsResponseDTO struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	ActivityValue float32              `json:"activity_value"`
	Exercises     []StudentExerciseDTO `json:"exercises"`
}
