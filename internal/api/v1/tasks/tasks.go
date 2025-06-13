package tasks

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/passwordhash/task-manager-api/internal/api/v1/response"
	"github.com/passwordhash/task-manager-api/internal/service"
)

type createTaskResponse struct {
	TaskUUID string `json:"task_uuid"`
}

func (h *handler) create(c *gin.Context) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	uuid, err := h.taskService.CreateTask(ctx)
	if response.HandleError(c, err) {
		return
	}

	response.NewOk(c, createTaskResponse{TaskUUID: uuid})
}

type statusResponse struct {
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Duration  string `json:"duration"`

	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (h *handler) status(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		response.NewErr(c, http.StatusBadRequest, response.ErrBadRequestParams, "Task UUID is required")
		return
	}

	task, err := h.taskService.Get(c, uuid)
	if errors.Is(err, service.ErrNotFound) {
		response.NewErr(c, http.StatusNotFound, response.ErrNotFound, "Task not found")
		return
	}
	if response.HandleError(c, err) {
		return
	}

	var taskErrResp string
	if task.Error != nil {
		taskErrResp = task.Error.Error()
	}

	response.NewOk(c, statusResponse{
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		Duration:  task.RunningDuration().String(),

		Result: task.Result,
		Error:  taskErrResp,
	})
}

type task struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type listTasksResponse struct {
	Tasks []task `json:"tasks"`
}

func (h *handler) list(c *gin.Context) {
	tasks, err := h.taskService.GetAll(c)
	if response.HandleError(c, err) {
		return
	}

	respTasks := make([]task, 0, len(tasks))
	for _, t := range tasks {
		respTasks = append(respTasks, task{
			UUID:   t.UUID,
			Status: string(t.Status),
		})
	}

	response.NewOk(c, listTasksResponse{Tasks: respTasks})
}

func (h *handler) cancel(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		response.NewErr(c, http.StatusBadRequest, response.ErrBadRequestParams, "Task UUID is required")
		return
	}

	err := h.taskService.Cancel(c, uuid)
	if errors.Is(err, service.ErrNotFound) {
		response.NewErr(c, http.StatusNotFound, response.ErrNotFound, "Task not found")
		return
	}
	if errors.Is(err, service.ErrCantCancel) {
		response.NewErr(c, http.StatusConflict, errors.New("cant_be_canceled"), "Task cannot be cancelled because it is already completed or cancelled")
		return
	}
	if response.HandleError(c, err) {
		return
	}

	response.NewOk(c, response.Message{Message: "Task cancelled successfully"})
}
