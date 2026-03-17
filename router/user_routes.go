package router

import (
	"schoolmanagement/internal/di"
	"schoolmanagement/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, container *di.Container) {
	user := r.Group("/user")
	{
		user.POST("/", container.UserHandler.CreateUser)
		user.GET("/", container.UserHandler.GetAllUsers)
		user.GET("/:id", container.UserHandler.GetUserByID)
		user.PATCH("/:id", middleware.AuthRoles("USER"), container.UserHandler.UpdateUser)
		user.DELETE("/:id", container.UserHandler.DeleteUser)
	}
}
