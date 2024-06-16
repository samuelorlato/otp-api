package handlers

import "github.com/gin-gonic/gin"

type HTTPErrorHandler struct{}

func NewHTTPErrorHandler() *HTTPErrorHandler {
	return &HTTPErrorHandler{}
}

// TODO: create custom error type
func (h *HTTPErrorHandler) Handle(err error, c *gin.Context) {
}
