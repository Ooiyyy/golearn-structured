package routes

import (
	"golearn-structured/internal/handler"

	"github.com/gin-gonic/gin"
)

func TodoRoute(r *gin.Engine, todoHandler *handler.TodoHandler) {

	api := r.Group("/api")
	{
		api.POST("/todos", todoHandler.Create)
		api.GET("/todos", todoHandler.GetAll)
		api.PATCH("/todos/:id", todoHandler.Update)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
	}
}
