package router

import (
	"schoolmanagement/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupActivityRoutes(r *gin.RouterGroup, container *di.Container) {
	activity := r.Group("/activity")
	{
		activity.POST("/", container.ActivityHandler.CreateActivity)
		activity.GET("/", container.ActivityHandler.GetAllActivities)
		activity.GET("/:id", container.ActivityHandler.GetActivityByID)
		activity.PATCH("/:id", container.ActivityHandler.UpdateActivity)
		activity.DELETE("/:id", container.ActivityHandler.DeleteActivity)
	}
}
