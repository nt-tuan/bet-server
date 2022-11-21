package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) login(c *gin.Context) {
	type LoginDto struct {
		Token string `json:"token"`
	}
	var loginDto LoginDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		s.responseError(c, err)
		return
	}

	fmt.Printf("loginDto %v\n", loginDto)

	username, err := s.authService.ParseIdToken(loginDto.Token)
	if err != nil {
		s.responseError(c, err)
		return
	}

	accessToken, err := s.authService.CreateAccessToken(username)
	if err != nil {
		s.responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}
