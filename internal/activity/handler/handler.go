package handler

import (
	"net/http"

	"schoolmanagement/internal/activity/dto"
	"schoolmanagement/internal/activity/service"
	"schoolmanagement/internal/utils"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	service service.ActivityService
}

func NewActivityHandler(service service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

// CreateActivity godoc
// @Summary Create a new activity
// @Description Create a new activity in the system.
// @Tags Activity
// @Accept json
// @Produce json
// @Param request body dto.ActivityRequestDTO true "Activity Info"
// @Success 200 {object} map[string]interface{} "Activity created successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/ [post]
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req dto.ActivityRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	activity, err := h.service.CreateActivity(&req)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade criada com sucesso.", activity)
}

// SubmitActivity godoc
// @Summary Submit an activity
// @Description Submit an activity response for a student.
// @Tags Activity
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param request body dto.SubmissionRequestDTO true "Submission Info"
// @Success 200 {object} map[string]interface{} "Activity submitted successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/submit/{id} [post]
func (h *ActivityHandler) SubmitActivity(c *gin.Context) {
	activityId := c.Param("id")
	userId := utils.GetUserId(c)

	var req dto.SubmissionRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	activity, err := h.service.SubmitActivity(&req, activityId, userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade submetida com sucesso.", activity)
}

// GetAllActivities godoc
// @Summary List all activities
// @Description Get a list of all activities.
// @Tags Activity
// @Produce json
// @Success 200 {object} map[string]interface{} "List of activities"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/ [get]
func (h *ActivityHandler) GetAllActivities(c *gin.Context) {
	activities, err := h.service.GetAllActivities()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Atividades listadas com sucesso.", activities)
}

// GetActivityByID godoc
// @Summary Get activity by ID
// @Description Find activity by their unique ID.
// @Tags Activity
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "Activity found"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /activity/{id} [get]
func (h *ActivityHandler) GetActivityByID(c *gin.Context) {
	id := c.Param("id")

	activity, err := h.service.GetActivityByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Atividade não encontrada.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Atividade encontrada.", activity)
}

// UpdateActivity godoc
// @Summary Update activity by ID
// @Description Partially update activity.
// @Tags Activity
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param request body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]interface{} "Activity updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/{id} [patch]
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	activity, err := h.service.UpdateActivity(id, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade atualizada com sucesso.", activity)
}

// DeleteActivity godoc
// @Summary Delete activity by ID
// @Description Delete activity from the system.
// @Tags Activity
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "Activity deleted successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/{id} [delete]
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteActivity(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade deletada com sucesso.", nil)
}

// GetActivityDashboard godoc
// @Summary Get activity dashboard
// @Description Get dashboard metrics for a specific activity.
// @Tags Activity
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "Dashboard data"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/{id}/dashboard [get]
func (h *ActivityHandler) GetActivityDashboard(c *gin.Context) {
	id := c.Param("id")

	dashboardData, err := h.service.GetActivityDashboard(id)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Dashboard da atividade gerado com sucesso.", dashboardData)
}

// ToggleActivityStatus godoc
// @Summary Toggle activity status
// @Description Release or hide an activity based on its current status.
// @Tags Activity
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "Activity status toggled"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /activity/{id}/release [post]
func (h *ActivityHandler) ToggleActivityStatus(c *gin.Context) {
	id := c.Param("id")

	// 1. Fetch current activity
	activity, err := h.service.GetActivityByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Atividade não encontrada.", http.StatusNotFound)
		return
	}

	// 2. Toggle status
	newStatus := "ACTIVE"
	if activity.Status == "ACTIVE" {
		newStatus = "INACTIVE"
	}

	// 3. Update status
	updates := map[string]interface{}{"status": newStatus}
	updatedActivity, err := h.service.UpdateActivity(id, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Status da atividade alterado.", updatedActivity)
}

// GetActiveActivities godoc
// @Summary Get active activities
// @Description Get active activities for the logged-in user.
// @Tags Activity
// @Produce json
// @Success 200 {object} map[string]interface{} "List of active activities"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/active [get]
func (h *ActivityHandler) GetActiveActivities(c *gin.Context) {
	userId := utils.GetUserId(c)

	activities, err := h.service.GetActiveActivities(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Atividades ativas listadas com sucesso.", activities)
}

// GetActivityQuestions godoc
// @Summary Get activity questions
// @Description Get questions for a specific activity.
// @Tags Activity
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "List of questions"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /activity/{id}/questions [get]
func (h *ActivityHandler) GetActivityQuestions(c *gin.Context) {
	id := c.Param("id")

	questions, err := h.service.GetActivityQuestions(id)
	if err != nil {
		utils.SendErrorResponse(c, "Atividade não encontrada.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Questões da atividade listadas com sucesso.", questions)
}

// GetStudentDashboard godoc
// @Summary Get student dashboard
// @Description Get dashboard insights for the logged-in student.
// @Tags Activity
// @Produce json
// @Success 200 {object} map[string]interface{} "Dashboard data"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/student/dashboard [get]
func (h *ActivityHandler) GetStudentDashboard(c *gin.Context) {
	userId := utils.GetUserId(c)

	dashboardData, err := h.service.GetStudentDashboard(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Dashboard do aluno gerado com sucesso.", dashboardData)
}

// UpdateAlternative godoc
// @Summary Update alternative
// @Description Update a specific alternative for an exercise.
// @Tags Activity
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param exerciseID path string true "Exercise ID"
// @Param alternativeID path string true "Alternative ID"
// @Param request body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]interface{} "Alternative updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/{id}/exercise/{exerciseID}/alternative/{alternativeID} [patch]
func (h *ActivityHandler) UpdateAlternative(c *gin.Context) {
	// Pega o ID da alternativa que vem na rota
	alternativeID := c.Param("alternativeID")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAlternative(alternativeID, updates); err != nil {
		utils.SendErrorResponse(c, "Erro ao atualizar alternativa: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Alternativa atualizada com sucesso.", nil)
}

// UpdateActivityFull godoc
// @Summary Update full activity
// @Description Replace full activity details.
// @Tags Activity
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param request body dto.ActivityRequestDTO true "Full Activity Data"
// @Success 200 {object} map[string]interface{} "Activity updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /activity/{id} [put]
func (h *ActivityHandler) UpdateActivityFull(c *gin.Context) {
	id := c.Param("id")

	var req dto.ActivityRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	activity, err := h.service.UpdateActivityFull(id, &req)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade atualizada com sucesso.", activity)
}

// GetClassRanking godoc
// @Summary Get class ranking
// @Description Get ranking of all students based on their grades.
// @Tags Activity
// @Produce json
// @Success 200 {object} map[string]interface{} "Ranking data"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/ranking [get]
func (h *ActivityHandler) GetClassRanking(c *gin.Context) {
	ranking, err := h.service.GetClassRanking()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Ranking da sala retornado com sucesso.", ranking)
}

// GetClassroomMetrics godoc
// @Summary Get classroom metrics
// @Description Get overall performance and completion metrics for the classroom.
// @Tags Activity
// @Produce json
// @Success 200 {object} map[string]interface{} "Metrics data"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /activity/metrics [get]
func (h *ActivityHandler) GetClassroomMetrics(c *gin.Context) {
	metrics, err := h.service.GetClassroomMetrics()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Métricas da sala retornadas com sucesso.", metrics)
}
