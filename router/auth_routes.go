// router/auth_routes.go
package router

import (
	"schoolmanagement/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, container *di.Container) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", container.AuthHandler.RegisterUser)
		auth.POST("/login", container.AuthHandler.LoginUser)
	}
}
