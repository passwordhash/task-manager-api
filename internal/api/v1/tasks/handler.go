package tasks

import (
	"context"

	"github.com/gin-gonic/gin"
)

type TaskRunner interface {
	Run(ctx context.Context)
}

type handler struct {
	taskRunner TaskRunner
}

func NewHandler(
	taskRunner TaskRunner,
) *handler {
	return &handler{
		taskRunner: taskRunner,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	tasksGroup := router.Group("/tasks")
	{
		tasksGroup.POST("/run", h.run)
	}
}
