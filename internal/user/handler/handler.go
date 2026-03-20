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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new student (USER) manually.
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.UserRequestDTO true "User Info"
// @Success 200 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /user/ [post]
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

// GetAllUsers godoc
// @Summary List all users
// @Description Get a list of all users.
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{} "List of students"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /user/ [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Alunos listados com sucesso.", users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Find user by their unique ID.
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User found"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.SendErrorResponse(c, "Aluno não encontrado.", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, "Aluno encontrado.", user)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete user from the system.
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(id); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Aluno deletado com sucesso.", nil)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Partially update user.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /user/{id} [patch]
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
