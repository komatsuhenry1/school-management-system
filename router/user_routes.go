package router

import (
	"schoolmanagement/internal/di"
	"schoolmanagement/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, container *di.Container) {
	user := r.Group("/user")
	{
		user.POST("/", middleware.AuthRoles("TEACHER"), container.UserHandler.CreateUser)
		user.GET("/", middleware.AuthRoles("TEACHER", "USER"), container.UserHandler.GetAllUsers)
		user.GET("/:id", middleware.AuthRoles("TEACHER", "USER"), container.UserHandler.GetUserByID)
		user.PATCH("/:id", middleware.AuthRoles("USER", "TEACHER"), container.UserHandler.UpdateUser)
		user.DELETE("/:id", middleware.AuthRoles("USER", "TEACHER"), container.UserHandler.DeleteUser)
	}
}
