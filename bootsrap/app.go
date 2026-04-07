package bootstrap

import (
	"golearn-structured/config"
	"golearn-structured/internal/handler"
	"golearn-structured/internal/repository"
	"golearn-structured/internal/routes"
	"golearn-structured/internal/service"
)

func Run() {
	cfg := config.LoadEnv()
	port := cfg.App.AppPort
	jwtSecret := cfg.App.JWTSecret
	db := config.ConnectDB(cfg.DB)

	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	userService := service.NewUserService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	authHandler := handler.NewAuthHandler(userService, jwtSecret)
	todoHandler := handler.NewTodoHandlers(todoService)

	r := routes.SetupRouter(
		authHandler,
		todoHandler,
		jwtSecret,
	)

	r.Run(":" + port)
}
