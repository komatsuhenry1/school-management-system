package router

import (
	"schoolmanagement/internal/di"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	container := di.NewContainer()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
			"http://localhost:3000",      // Para acesso local via localhost
			"http://192.168.18.153:3000", // Para acesso via IP na rede local
			"http://192.168.18.190:3000", // Para acesso via IP na rede local
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")

	SetupAuthRoutes(api, container)
	SetupUserRoutes(api, container)

	return router
}
