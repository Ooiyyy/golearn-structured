package routes

import (
	"golearn-structured/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, todoHandler *handler.TodoHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	AuthRoute(r, authHandler, jwtSecret)
	TodoRoute(r, todoHandler)

	return r
}
