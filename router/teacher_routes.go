package router

import (
	"schoolmanagement/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupTeacherRoutes(r *gin.RouterGroup, container *di.Container) {
	teacher := r.Group("/teacher")
	{
		teacher.POST("/", container.TeacherHandler.CreateTeacher)
		teacher.GET("/", container.TeacherHandler.GetAllTeachers)
		teacher.GET("/:id", container.TeacherHandler.GetTeacherByID)
		teacher.PATCH("/:id", container.TeacherHandler.UpdateTeacher)
		teacher.DELETE("/:id", container.TeacherHandler.DeleteTeacher)
	}
}
