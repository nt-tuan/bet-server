package bet

import (
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotAbleChangeDealStatus error = errors.New("not-able-change-deal-status")
)

func (s *BetService) getDealCollection() *mongo.Collection {
	log.Print(s.database)
	return s.client.Database(s.database).Collection("deal")
}

func (s *BetService) CountDeals() (int64, error) {
	var count, err = s.getDealCollection().CountDocuments(s.ctx, bson.D{{}})
	return count, err
}

func (s *BetService) GetDeals() ([]Deal, error) {
	s.CountDeals()
	var cursor, err = s.getDealCollection().Find(s.ctx, bson.D{{}})

	if err == mongo.ErrNilDocument {
		return []Deal{}, nil
	}
	if err != nil {
		return nil, err
	}

	deals := []Deal{}
	if err := cursor.All(s.ctx, &deals); err != nil {
		return nil, err
	}
	return deals, nil
}

func (s *BetService) GetDeal(id string) (*Deal, error) {
	var objectId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	var result = s.getDealCollection().FindOne(s.ctx, bson.D{{
		Key:   "_id",
		Value: objectId,
	}})

	if result == nil {
		return nil, mongo.ErrNoDocuments
	}

	var deal Deal
	result.Decode(&deal)
	return &deal, nil
}

func (s *BetService) CreateDeal(title string, infoSource string, options []struct {
	Title      string
	InfoSource string
}, startDate time.Time, closedDate time.Time, createdBy string) error {
	var dealOptions = []DealOption{}
	for _, option := range options {
		dealOptions = append(dealOptions, DealOption{
			ID:         primitive.NewObjectID(),
			Title:      option.Title,
			InfoSource: option.InfoSource,
		})
	}
	var deal = Deal{
		ID:         primitive.NewObjectID(),
		Title:      title,
		InfoSource: infoSource,
		Status:     DealStatusInactive,
		Options:    dealOptions,
		StartDate:  startDate,
		ClosedDate: closedDate,
		CreatedBy:  createdBy,
	}
	var collection = s.getDealCollection()
	var _, err = collection.InsertOne(s.ctx, deal)
	return err
}

func (s *BetService) updateDealStatus(id string, status DealStatus) error {
	var deal, getDealErr = s.GetDeal(id)
	if getDealErr != nil {
		return getDealErr
	}

	if !deal.canChangeStatus() {
		return ErrNotAbleChangeDealStatus
	}
	var objectId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "status", Value: status}},
	}}
	var _, updateErr = s.getDealCollection().UpdateByID(s.ctx, objectId, update)
	return updateErr
}

func (s *BetService) OpenDeal(id string) error {
	return s.updateDealStatus(id, DealStatusActive)

}

func (s *BetService) CancelDeal(id string) error {
	return s.updateDealStatus(id, DealStatusCanceled)
}

func checkDealOptionID(deal Deal, optionID string) error {
	for _, s := range deal.Options {
		if s.ID.Hex() == optionID {
			return nil
		}
	}
	return ErrInvalidDealOption
}
