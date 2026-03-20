package handler

import (
	"schoolmanagement/internal/auth/service"

	"github.com/gin-gonic/gin"
	"schoolmanagement/internal/auth/dto"
	"schoolmanagement/internal/utils"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user in the system with role validation.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.UserRequestDTO true "User Registration Info"
// @Success 200 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /auth/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var userRequestDto dto.UserRequestDTO

	if err := c.ShouldBindJSON(&userRequestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userRequestDto.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RegisterUser(&userRequestDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário registrado com sucesso.", nil)
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequestDTO true "User Login Info"
// @Success 200 {object} map[string]interface{} "User logged in successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /auth/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginRequestDTO dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginRequestDTO); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	token, refreshToken, user, err := h.service.LoginUser(&loginRequestDTO)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	c.SetCookie(
		"refresh_token", // nome
		refreshToken,    // valor
		60*60*24*7,      // duração em segundos (7 dias)
		"/",             // path
		"",              // domain (coloque seu domínio em prod)
		false,           // secure (true em produção HTTPS)
		true,            // httpOnly (IMPORTE!)
	)

	utils.SendSuccessResponse(c, "Usuário logado com sucesso.",
		gin.H{
			"token": token,
			"user":  gin.H{"name": user.Name, "email": user.Email, "role": user.Role, "id": user.ID},
		})
}
