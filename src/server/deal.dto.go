package server

import "time"

type BaseInfo struct {
	Title      string `json:"title" binding:"required"`
	InfoSource string `json:"infoSource" binding:"required"`
}
type DealDto struct {
	*BaseInfo
	Options    []BaseInfo `json:"options" binding:"required"`
	StartDate  time.Time  `json:"startDate" binding:"required"`
	ClosedDate time.Time  `json:"closedDate" binding:"required"`
}
