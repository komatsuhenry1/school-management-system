package di

import (
	"schoolmanagement/config"
	activityHandler "schoolmanagement/internal/activity/handler"
	activityRepository "schoolmanagement/internal/activity/repository"
	activityService "schoolmanagement/internal/activity/service"
	authHandler "schoolmanagement/internal/auth/handler"
	authService "schoolmanagement/internal/auth/service"
	teacherHandler "schoolmanagement/internal/teacher/handler"
	teacherRepository "schoolmanagement/internal/teacher/repository"
	teacherService "schoolmanagement/internal/teacher/service"
	userHandler "schoolmanagement/internal/user/handler"
	userRepository "schoolmanagement/internal/user/repository"
	userService "schoolmanagement/internal/user/service"
)

type Container struct {
	AuthHandler     *authHandler.UserHandler
	UserHandler     *userHandler.UserHandler
	ActivityHandler *activityHandler.ActivityHandler
	TeacherHandler  *teacherHandler.TeacherHandler
}

func NewContainer() *Container {
	db := config.GetDB()

	// User (shared repository)
	userRepo := userRepository.NewUserRepository(db)

	// Auth
	authSvc := authService.NewUserService(userRepo)
	authHdl := authHandler.NewUserHandler(authSvc)

	// User
	userSvc := userService.NewUserService(userRepo)
	userHdl := userHandler.NewUserHandler(userSvc)

	// Activity
	activityRepo := activityRepository.NewActivityRepository(db)
	activitySvc := activityService.NewActivityService(activityRepo)
	activityHdl := activityHandler.NewActivityHandler(activitySvc)

	// Teacher
	teacherRepo := teacherRepository.NewTeacherRepository(db)
	teacherSvc := teacherService.NewTeacherService(teacherRepo)
	teacherHdl := teacherHandler.NewTeacherHandler(teacherSvc)

	return &Container{
		AuthHandler:     authHdl,
		UserHandler:     userHdl,
		ActivityHandler: activityHdl,
		TeacherHandler:  teacherHdl,
	}
}
