package server

import (
	"bet-server/src/bet"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getDealById(c *gin.Context) (*bet.Deal, error) {
	var id, _ = c.Params.Get("id")
	var deal, err = s.betService.GetDeal(id)
	if err != nil {
		s.responseError(c, err)
		return nil, err
	}

	return deal, nil
}

func (s *Server) getDealController(c *gin.Context) {
	var deal, err = s.getDealById(c)
	if err != nil {
		s.responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (s *Server) getDeals(c *gin.Context) {
	var deals, err = s.betService.GetDeals()
	if err != nil {
		s.responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, deals)
}

func (s *Server) createDeal(c *gin.Context) {
	var deal DealDto
	if err := c.ShouldBindJSON(&deal); err != nil {
		s.responseError(c, err)
		return
	}

	type Option struct {
		Title      string
		InfoSource string
	}
	var options = []struct {
		Title      string
		InfoSource string
	}{}
	for _, option := range deal.Options {
		options = append(options, struct {
			Title      string
			InfoSource string
		}{
			Title:      option.Title,
			InfoSource: option.InfoSource,
		})
	}

	s.betService.CreateDeal(deal.Title, deal.InfoSource, options, deal.StartDate, deal.ClosedDate, s.getUser(c))
	c.Status(http.StatusCreated)
}

func (s *Server) openDeal(c *gin.Context) {
	var bet, err = s.getDealById(c)
	if err != nil {
		s.responseError(c, err)
		return
	}

	var openErr = s.betService.OpenDeal(bet.ID.Hex())
	if openErr != nil {
		s.responseError(c, openErr)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) cancelDeal(c *gin.Context) {
	var bet, err = s.getDealById(c)
	if err != nil {
		return
	}

	var openErr = s.betService.CancelDeal(bet.ID.Hex())
	if openErr != nil {
		s.responseError(c, err)
	}

	c.Status(http.StatusOK)
}

func (s *Server) getHighlightDeal(c *gin.Context) {
	var deal, err = s.betService.GetHighlightDeal()
	if err != nil {
		s.responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, deal)
}
