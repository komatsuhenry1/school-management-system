package main

import (
	"schoolmanagement/config"
	"schoolmanagement/router"
	"os"
)

// @title School Management System API
// @version 1.0
// @description API for the School Management System with roles for Teachers and Students.
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// go middleware.CleanupClients()
	if err := config.ConnectDB(); err != nil {
    	panic(err) // ou log.Fatal(err)
  	}
	// utils.CreateKeys()
	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run("0.0.0.0:" + port)
}
