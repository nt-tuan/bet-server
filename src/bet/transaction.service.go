package bet

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *BetService) getTransactionCollection() *mongo.Collection {
	return s.client.Database(s.database).Collection("transaction")
}

func (s *BetService) createTransaction(txType TransactionType, amount int, source string, reference string, meta *bson.M) error {
	var _, err = s.getTransactionCollection().InsertOne(s.ctx, Transaction{
		ID:        primitive.NewObjectID(),
		Type:      txType,
		Amount:    amount,
		Source:    source,
		Reference: reference,
		Meta:      meta,
	})
	return err
}
