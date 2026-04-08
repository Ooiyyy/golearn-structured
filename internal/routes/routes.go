package routes

import (
	"golearn-structured/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, todoHandler *handler.TodoHandler, jwtSecret string) *gin.Engine {
	// gin.Default() menambahkan logger + recovery middleware bawaan.
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	// Daftarkan endpoint per modul.
	AuthRoute(r, authHandler, jwtSecret)
	TodoRoute(r, todoHandler)

	// Router ini dipanggil oleh bootstrap untuk menerima request dari client.
	return r
}
