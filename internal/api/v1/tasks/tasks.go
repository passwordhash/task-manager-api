package tasks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) run(c *gin.Context) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	h.taskRunner.Run(ctx)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
