package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) responseError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (s *Server) getUser(c *gin.Context) string {
	var user, ok = c.Get("user")
	if !ok {
		return "unknown"
	}
	return user.(string)
}
