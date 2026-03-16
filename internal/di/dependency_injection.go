package di

import (
	"schoolmanagement/config"
	authHandler "schoolmanagement/internal/auth/handler"
	authService "schoolmanagement/internal/auth/service"
	userHandler "schoolmanagement/internal/user/handler"
	userRepository "schoolmanagement/internal/user/repository"
	userService "schoolmanagement/internal/user/service"
)

type Container struct {
	AuthHandler        *authHandler.UserHandler
	UserHandler        *userHandler.UserHandler
}

func NewContainer() *Container {
	db := config.GetDB()
	//supabaseClient := config.GetClient()

	userRepo := userRepository.NewUserRepository(db)

	authSvc := authService.NewUserService(userRepo)
	authHdl := authHandler.NewUserHandler(authSvc)

	userSvc := userService.NewUserService(userRepo)
	userHdl := userHandler.NewUserHandler(userSvc)

	return &Container{
		AuthHandler:        authHdl,
		UserHandler:        userHdl,
	}
}
