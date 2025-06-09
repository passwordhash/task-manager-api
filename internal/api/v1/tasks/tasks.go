package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) run(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
