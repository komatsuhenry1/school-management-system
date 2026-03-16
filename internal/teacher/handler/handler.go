package handler

import (
	"net/http"

	"schoolmanagement/internal/teacher/dto"
	"schoolmanagement/internal/teacher/service"
	"schoolmanagement/internal/utils"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	service service.TeacherService
}

func NewTeacherHandler(service service.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: service}
}

func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req dto.TeacherRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	teacher, err := h.service.CreateTeacher(&req)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Professor criado com sucesso.", teacher)
}

func (h *TeacherHandler) GetAllTeachers(c *gin.Context) {
	teachers, err := h.service.GetAllTeachers()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Professores listados com sucesso.", teachers)
}

func (h *TeacherHandler) GetTeacherByID(c *gin.Context) {
	id := c.Param("id")

	teacher, err := h.service.GetTeacherByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Professor não encontrado.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Professor encontrado.", teacher)
}

func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	teacher, err := h.service.UpdateTeacher(id, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Professor atualizado com sucesso.", teacher)
}

func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteTeacher(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Professor deletado com sucesso.", nil)
}
