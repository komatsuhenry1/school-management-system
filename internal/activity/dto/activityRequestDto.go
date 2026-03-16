package dto

type ActivityRequestDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
