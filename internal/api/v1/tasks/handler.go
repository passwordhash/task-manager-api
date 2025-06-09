package tasks

import "github.com/gin-gonic/gin"

type

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	tasksGroup := router.Group("/tasks")
	{
		tasksGroup.POST("/run", h.run)
	}
}
