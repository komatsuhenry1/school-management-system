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
		activity.PUT("/:id", middleware.AuthRoles("TEACHER"), container.ActivityHandler.UpdateActivityFull)
		activity.PATCH("/:id", middleware.AuthRoles("TEACHER"), container.ActivityHandler.UpdateActivity)
		activity.PATCH("/:id/exercise/:exerciseID/alternative/:alternativeID", middleware.AuthRoles("TEACHER"), container.ActivityHandler.UpdateAlternative)
		activity.DELETE("/:id", middleware.AuthRoles("TEACHER"), container.ActivityHandler.DeleteActivity)
		
		// Somente aluno responde
		activity.POST("/submit/:id", middleware.AuthRoles("USER"), container.ActivityHandler.SubmitActivity)
		activity.GET("/student/dashboard", middleware.AuthRoles("USER"), container.ActivityHandler.GetStudentDashboard)
		
		// Ranking e Metrics (Ambos teachers and students can see)
		activity.GET("/ranking", middleware.AuthRoles("USER", "TEACHER"), container.ActivityHandler.GetClassRanking)
		activity.GET("/metrics", middleware.AuthRoles("USER", "TEACHER"), container.ActivityHandler.GetClassroomMetrics)

		// Ambos podem listar e ver atividades específicas (Teacher vê tudo, Activity/active aluno vê)
		activity.GET("/", middleware.AuthRoles("TEACHER", "USER"), container.ActivityHandler.GetAllActivities)
		activity.GET("/active", middleware.AuthRoles("TEACHER", "USER"), container.ActivityHandler.GetActiveActivities)
		activity.GET("/:id", middleware.AuthRoles("TEACHER", "USER"), 	container.ActivityHandler.GetActivityByID)
		activity.GET("/:id/questions", middleware.AuthRoles("TEACHER", "USER"), container.ActivityHandler.GetActivityQuestions)
		
		// Dashboard apenas professores
		activity.GET("/:id/dashboard", middleware.AuthRoles("TEACHER"), container.ActivityHandler.GetActivityDashboard)
		
		// Liberar/Ocultar atividade
		activity.POST("/:id/release", middleware.AuthRoles("TEACHER"), container.ActivityHandler.ToggleActivityStatus)
	}
}
