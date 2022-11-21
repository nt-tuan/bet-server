package bet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BetStatus = string

const (
	BetStatusPending  BetStatus = "PENDING"
	BetStatusActive   BetStatus = "ACTIVE"
	BetStatusCanceled BetStatus = "CANCELED"
)

type Bet struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	DealId       primitive.ObjectID `bson:"dealID" json:"dealID"`
	DealOptionId primitive.ObjectID `bson:"dealOptionID" json:"dealOptionID"`

	Amount    int       `bson:"amount" json:"amount"`
	Status    BetStatus `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
}
