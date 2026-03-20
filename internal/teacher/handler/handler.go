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

// CreateTeacher godoc
// @Summary Create a new teacher
// @Description Create a new teacher manually.
// @Tags Teacher
// @Accept json
// @Produce json
// @Param request body dto.TeacherRequestDTO true "Teacher Info"
// @Success 200 {object} map[string]interface{} "Teacher created successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /teacher/ [post]
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

// GetAllTeachers godoc
// @Summary List all teachers
// @Description Get a list of all teachers.
// @Tags Teacher
// @Produce json
// @Success 200 {object} map[string]interface{} "List of teachers"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /teacher/ [get]
func (h *TeacherHandler) GetAllTeachers(c *gin.Context) {
	teachers, err := h.service.GetAllTeachers()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Professores listados com sucesso.", teachers)
}

// GetTeacherByID godoc
// @Summary Get teacher by ID
// @Description Find teacher by their unique ID.
// @Tags Teacher
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {object} map[string]interface{} "Teacher found"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /teacher/{id} [get]
func (h *TeacherHandler) GetTeacherByID(c *gin.Context) {
	id := c.Param("id")

	teacher, err := h.service.GetTeacherByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Professor não encontrado.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Professor encontrado.", teacher)
}

// UpdateTeacher godoc
// @Summary Update teacher by ID
// @Description Partially update teacher.
// @Tags Teacher
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param request body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]interface{} "Teacher updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /teacher/{id} [patch]
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

// DeleteTeacher godoc
// @Summary Delete teacher by ID
// @Description Delete teacher from the system.
// @Tags Teacher
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {object} map[string]interface{} "Teacher deleted successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /teacher/{id} [delete]
func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteTeacher(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Professor deletado com sucesso.", nil)
}
