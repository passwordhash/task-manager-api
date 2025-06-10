package tasks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTaskResponse struct {
	TaskUUID string `json:"task_uuid"`
}

func (h *handler) create(c *gin.Context) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	uuid, err := h.taskService.CreateTask(ctx)
	// TODO: handle error properly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
		})
		return
	}

	c.JSON(http.StatusOK, createTaskResponse{
		TaskUUID: uuid,
	})
}
