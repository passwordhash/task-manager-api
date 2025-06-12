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

type task struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type listTasksResponse struct {
	Tasks []task `json:"tasks"`
}

func (h *handler) list(c *gin.Context) {
	tasks, err := h.taskService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list tasks",
		})
		return
	}

	respTasks := make([]task, 0, len(tasks))
	for _, t := range tasks {
		respTasks = append(respTasks, task{
			UUID:   t.UUID,
			Status: string(t.Status),
		})
	}

	c.JSON(http.StatusOK, listTasksResponse{
		Tasks: respTasks,
	})
}

func (h *handler) cancel(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task UUID is required",
		})
		return
	}

	err := h.taskService.Cancel(c, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cancel task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task cancelled successfully",
	})
}
