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
	var req dto.SubmissionRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	activity, err := h.service.SubmitActivity(&req)
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
