package tasks

import (
	"context"

	"github.com/gin-gonic/gin"
)

type TaskManager interface {
	Run(ctx context.Context)
}

type handler struct {
	taskManager TaskManager
}

func NewHandler(
	taskRunner TaskManager,
) *handler {
	return &handler{
		taskManager: taskRunner,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	tasksGroup := router.Group("/tasks")
	{
		tasksGroup.POST("/run", h.run)
	}
}
