package bootstrap

import (
	"golearn-structured/config"
	"golearn-structured/internal/handler"
	"golearn-structured/internal/repository"
	"golearn-structured/internal/routes"
	"golearn-structured/internal/service"
)

func Run() {
	// 1) Load konfigurasi runtime (port, JWT secret, kredensial DB) dari env.
	cfg := config.LoadEnv()
	port := cfg.App.AppPort
	jwtSecret := cfg.App.JWTSecret
	// 2) Buka koneksi DB pool sekali saat startup, dipakai bersama di semua request.
	db := config.ConnectDB(cfg.DB)

	// 3) Rangkai dependensi dari lapisan paling bawah ke atas:
	// repository -> service -> handler.
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

	// 4) Mulai HTTP server. Setelah ini request masuk via router Gin.
	r.Run(":" + port)
}
