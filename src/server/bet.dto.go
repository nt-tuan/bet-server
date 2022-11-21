package server

type CreateBetDto struct {
	DealID       string `json:"dealID"`
	DealOptionID string `json:"dealOptionID"`
	Amount       int    `json:"amount" binding:"required,gte=1,lte=100"`
}
