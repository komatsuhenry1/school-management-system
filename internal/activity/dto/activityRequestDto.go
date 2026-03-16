package dto

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
	ActivityValue float32       `json:"activity_value"`
	Exercises     []ExerciseDTO `json:"exercises"`
}
