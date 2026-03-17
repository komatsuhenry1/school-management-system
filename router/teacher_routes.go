package router

import (
	"schoolmanagement/internal/di"
	"schoolmanagement/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupTeacherRoutes(r *gin.RouterGroup, container *di.Container) {
	teacher := r.Group("/teacher")
	{
		teacher.POST("/", middleware.AuthRoles("TEACHER"), container.TeacherHandler.CreateTeacher)
		teacher.GET("/", middleware.AuthRoles("TEACHER"), container.TeacherHandler.GetAllTeachers)
		teacher.GET("/:id", middleware.AuthRoles("TEACHER"), container.TeacherHandler.GetTeacherByID)
		teacher.PATCH("/:id", middleware.AuthRoles("TEACHER"), container.TeacherHandler.UpdateTeacher)
		teacher.DELETE("/:id", middleware.AuthRoles("TEACHER"), container.TeacherHandler.DeleteTeacher)
	}
}
