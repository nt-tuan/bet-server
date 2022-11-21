package bet

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	TransactionTypeCreateBet = "CREATE_BET"
	TransactionTypeCancelBet = "CANCEL_BET"
	TransactionTypeWithdraw  = "WITHDRAW"
	TransactionTypeDeposit   = "DEPOSIT"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      TransactionType    `bson:"type" json:"type"`
	Amount    int                `bson:"amount" json:"amount"`
	Source    string             `bson:"source" json:"source"`
	Reference string             `bson:"reference" json:"reference"`
	Meta      *bson.M            `bson:"meta" json:"meta"`
}
