package router

import (
	"schoolmanagement/internal/di"
	"schoolmanagement/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupActivityRoutes(r *gin.RouterGroup, container *di.Container) {
	activity := r.Group("/activity")
	{
		// Somente professor cria, atualiza e deleta
		activity.POST("/", middleware.AuthRoles("TEACHER"), container.ActivityHandler.CreateActivity)
		activity.PATCH("/:id", middleware.AuthRoles("TEACHER"), container.ActivityHandler.UpdateActivity)
		activity.DELETE("/:id", middleware.AuthRoles("TEACHER"), container.ActivityHandler.DeleteActivity)
		
		// Somente aluno responde
		activity.POST("/submit", middleware.AuthRoles("USER"), container.ActivityHandler.SubmitActivity)
		
		// Ambos podem listar e ver atividades específicas
		activity.GET("/", middleware.AuthRoles("TEACHER", "USER"), container.ActivityHandler.GetAllActivities)
		activity.GET("/:id", middleware.AuthRoles("TEACHER", "USER"), container.ActivityHandler.GetActivityByID)
		
		// Dashboard apenas professores
		activity.GET("/:id/dashboard", middleware.AuthRoles("TEACHER"), container.ActivityHandler.GetActivityDashboard)
		
		// Liberar/Ocultar atividade
		activity.POST("/:id/release", middleware.AuthRoles("TEACHER"), container.ActivityHandler.ToggleActivityStatus)
	}
}
