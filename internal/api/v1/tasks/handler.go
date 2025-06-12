package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/passwordhash/task-manager-api/internal/service"
)

type handler struct {
	taskService service.TaskService
}

func NewHandler(
	taskService service.TaskService,
) *handler {
	return &handler{
		taskService: taskService,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	tasksGroup := router.Group("/tasks")
	{
		tasksGroup.GET("/", h.list)
		tasksGroup.POST("/create", h.create)

		taskGroup := tasksGroup.Group("/:uuid")
		{
			taskGroup.POST("/cancel", h.cancel)
		}
	}
}
