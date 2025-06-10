package tasks

import (
	"context"

	"github.com/gin-gonic/gin"
)

type TaskManager interface {
	CreateTask(ctx context.Context) (taskUUID string, err error)
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
		tasksGroup.POST("/create", h.create)
	}
}
