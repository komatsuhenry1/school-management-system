package dto

type TeacherRequestDTO struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
}
