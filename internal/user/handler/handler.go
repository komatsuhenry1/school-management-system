package handler

import (
	"fmt"
	"net/http"
	"schoolmanagement/internal/auth/dto"
	"schoolmanagement/internal/user/service"
	"schoolmanagement/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service     service.UserService
	authService interface {
		RegisterUser(userRequestDTO *dto.UserRequestDTO) error
	}
}

// We need to inject the authService to use its RegisterUser function since
// the user creation logic (password hashing, valid cpf/phone) is in the auth package.
func NewUserHandler(service service.UserService, authService interface {
	RegisterUser(userRequestDTO *dto.UserRequestDTO) error
}) *UserHandler {
	return &UserHandler{
		service:     service,
		authService: authService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequestDto dto.UserRequestDTO

	if err := c.ShouldBindJSON(&userRequestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userRequestDto.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.RegisterUser(&userRequestDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário registrado com sucesso.", nil)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Alunos listados com sucesso.", users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Aluno não encontrado.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Aluno encontrado.", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Aluno deletado com sucesso.", nil)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	protectedFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"password":   true, // Extra safety for user entity updates
	}

	for key := range updates {
		if protectedFields[strings.ToLower(key)] {
			utils.SendErrorResponse(c, fmt.Sprintf("Campo(s) %s não pode ser atualizado.", key), http.StatusBadRequest)
			return
		}
	}

	user, err := h.service.UpdateUserPartial(userId, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário atualizado com sucesso.", user)
}
