package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequestParams = errors.New("invalid_request_parameters")
	ErrNotFound         = errors.New("not_found")
)

type Message struct {
	Message string `json:"message"`
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func NewErr(c *gin.Context, code int, err error, clientMessage string) {
	c.JSON(code, Error{
		Error:   err.Error(),
		Message: clientMessage,
	})
}

// HandleError processes handlgin basic errors in the context of a gin handler.
// It returns true if an error was handled, false otherwise.
func HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		NewErr(c, http.StatusGatewayTimeout, errors.New("gateway_timeout"), "The request timed out.")
		return true
	}
	if errors.Is(err, context.Canceled) {
		NewErr(c, http.StatusRequestTimeout, errors.New("request_timeout"), "The request was canceled.")
		return true
	}
	NewErr(c, http.StatusInternalServerError, err, "Unexpected error occurred.")
	return true
}
