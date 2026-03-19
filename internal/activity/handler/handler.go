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

func (h *ActivityHandler) GetAllActivities(c *gin.Context) {
	activities, err := h.service.GetAllActivities()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Atividades listadas com sucesso.", activities)
}

func (h *ActivityHandler) GetActivityByID(c *gin.Context) {
	id := c.Param("id")

	activity, err := h.service.GetActivityByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Atividade não encontrada.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Atividade encontrada.", activity)
}

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

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteActivity(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Atividade deletada com sucesso.", nil)
}

func (h *ActivityHandler) GetActivityDashboard(c *gin.Context) {
	id := c.Param("id")

	dashboardData, err := h.service.GetActivityDashboard(id)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Dashboard da atividade gerado com sucesso.", dashboardData)
}

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

func (h *ActivityHandler) GetActiveActivities(c *gin.Context) {
	userId := utils.GetUserId(c)

	activities, err := h.service.GetActiveActivities(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Atividades ativas listadas com sucesso.", activities)
}

func (h *ActivityHandler) GetActivityQuestions(c *gin.Context) {
	id := c.Param("id")

	questions, err := h.service.GetActivityQuestions(id)
	if err != nil {
		utils.SendErrorResponse(c, "Atividade não encontrada.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Questões da atividade listadas com sucesso.", questions)
}

func (h *ActivityHandler) GetStudentDashboard(c *gin.Context) {
	userId := utils.GetUserId(c)

	dashboardData, err := h.service.GetStudentDashboard(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Dashboard do aluno gerado com sucesso.", dashboardData)
}

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

func (h *ActivityHandler) GetClassRanking(c *gin.Context) {
	ranking, err := h.service.GetClassRanking()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Ranking da sala retornado com sucesso.", ranking)
}

func (h *ActivityHandler) GetClassroomMetrics(c *gin.Context) {
	metrics, err := h.service.GetClassroomMetrics()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Métricas da sala retornadas com sucesso.", metrics)
}
