package routes

import (
	"golearn-structured/internal/handler"
	"golearn-structured/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine, authHandler *handler.AuthHandler, jwtSecret string) {

	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)

	protected := r.Group("/")
	protected.Use(middleware.JWT(jwtSecret))
	{
		protected.GET("/profile", authHandler.Profile)
		protected.POST("/logout", authHandler.Logout)
	}
}
