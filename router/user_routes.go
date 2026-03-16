package router

import (
	"schoolmanagement/internal/di"
	"schoolmanagement/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, container *di.Container) {
	user := r.Group("/user")
	{
		user.PATCH("/:id", middleware.AuthUser(), container.UserHandler.UpdateUser)
	}
}
