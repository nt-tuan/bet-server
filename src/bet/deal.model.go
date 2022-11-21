package bet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealStatus = string

const (
	DealStatusInactive DealStatus = "INACTIVE"
	DealStatusActive   DealStatus = "ACTIVE"

	// Immutable Statuses
	DealStatusClosed   DealStatus = "CLOSED"
	DealStatusCanceled DealStatus = "CANCELED"
)

type DealOption struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Title      string             `bson:"title" json:"title"`
	InfoSource string             `bson:"infoSource" json:"infoSource"`
}

type Deal struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Title      string             `bson:"title" json:"title"`
	InfoSource string             `bson:"infoSource" json:"infoSource"`
	Status     DealStatus         `bson:"status" json:"status"`
	Options    []DealOption       `bson:"options" json:"options"`
	StartDate  time.Time          `bson:"startDate" json:"startDate"`
	ClosedDate time.Time          `bson:"closedDate" json:"closeDate"`
	CreatedBy  string             `bson:"createdBy" json:"createdBy"`
}

func (deal *Deal) isActive() bool {
	return deal.Status == DealStatusActive && deal.ClosedDate.After(time.Now())
}

var immutatableStatuses = []DealStatus{DealStatusClosed, DealStatusCanceled}

func (deal *Deal) canChangeStatus() bool {
	for _, status := range immutatableStatuses {
		if deal.Status == status {
			return false
		}
	}
	return true
}
