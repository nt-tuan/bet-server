package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) placeBet(c *gin.Context) {
	var bet CreateBetDto
	if err := c.ShouldBindJSON(&bet); err != nil {
		s.responseError(c, err)
		return
	}
	var err = s.betService.PlaceBet(bet.DealID, bet.DealOptionID, bet.Amount, s.getUser(c))
	if err != nil {
		s.responseError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func (s *Server) getBets(c *gin.Context) {
	var dealId, _ = c.Params.Get("id")

	var bets, err = s.betService.GetBets(dealId)
	if err != nil {
		s.responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bets)
}
