package bet

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDealNotInactive     = errors.New("deal-not-inactive")
	ErrDealStatusImmutated = errors.New("deal-status-immutated")
	ErrDealNotActive       = errors.New("deal-not-active")
	ErrInvalidID           = errors.New("invalid-id")
	ErrInvalidDealOption   = errors.New("invalid-deal-option")
	ErrBetNotFound         = errors.New("bet-not-found")
	ErrBetInactive         = errors.New("bet-inactive")
)

type BetService struct {
	client   *mongo.Client
	database string
	ctx      context.Context
}

var collection *mongo.Collection

func NewBetService(uri string, database string) BetService {
	// serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return BetService{client: client, database: database, ctx: ctx}
}

func (s *BetService) getBetCollection() *mongo.Collection {
	return s.client.Database(s.database).Collection("bet")
}

func (s *BetService) updateBetStatus(id string, betStatus BetStatus) error {
	var objectId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var _, updateErr = s.getBetCollection().UpdateOne(s.ctx, bson.D{{
		Key:   "_id",
		Value: objectId,
	}}, bson.D{{
		Key: "status", Value: betStatus,
	}})

	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (s *BetService) PlaceBet(dealId string, dealOptionId string, amount int, createdBy string) error {
	var deal, err = s.GetDeal(dealId)
	if err != nil {
		return err
	}

	if deal.Status != DealStatusActive {
		return errors.New("deal-no-active")
	}

	var now = time.Now()

	if now.Before(deal.StartDate) {
		return errors.New("deal-not-start-yet")
	}

	if now.After(deal.ClosedDate) {
		return errors.New("deal-is-closed")
	}

	dealObjectId, err := primitive.ObjectIDFromHex(dealId)
	if err != nil {
		return ErrInvalidID
	}

	if err := checkDealOptionID(*deal, dealOptionId); err != nil {
		return err
	}
	dealOptionObjectId, err := primitive.ObjectIDFromHex(dealOptionId)
	if err != nil {
		return ErrInvalidID
	}

	var bet = Bet{
		ID:           primitive.NewObjectID(),
		DealId:       dealObjectId,
		DealOptionId: dealOptionObjectId,
		Amount:       amount,
		Status:       BetStatusPending,
		CreatedAt:    time.Now(),
		CreatedBy:    createdBy,
	}

	log.Printf("%v", bet)

	var _, insertErr = s.getBetCollection().InsertOne(s.ctx, bet)
	return insertErr
}

func (s *BetService) checkBetTransaction(bet Bet) error {
	var deal, err = s.GetDeal(bet.DealId.Hex())
	if err != nil {
		return err
	}

	// Bet can not be inactive
	if bet.Status != BetStatusActive {
		return ErrBetInactive
	}

	// Deal must be active
	if deal.isActive() {
		return ErrDealNotActive
	}
	return nil
}

func (s *BetService) useBetTransaction(betId string, sess func(mongo.SessionContext, mongo.Session, Bet) error) error {
	var collection = s.getBetCollection()
	var result = collection.FindOne(s.ctx, bson.D{{
		Key:   "_id",
		Value: "",
	}})

	var bet Bet
	if err := result.Decode(&bet); err != nil {
		return err
	}

	if s.checkBetTransaction(bet) != nil {
		return ErrDealNotActive
	}

	if result == nil {
		return ErrBetNotFound
	}

	session, err := s.client.StartSession()
	if err != nil {
		return err
	}

	if err = session.StartTransaction(); err != nil {
		return err
	}
	sessionErr := mongo.WithSession(s.ctx, session, func(context mongo.SessionContext) error {
		return sess(context, session, bet)
	})
	if sessionErr != nil {
		session.EndSession(s.ctx)
		return err
	}

	session.EndSession(s.ctx)
	return nil
}

func (s *BetService) CancelBet(betId string) error {
	return s.updateBetStatus(betId, BetStatusCanceled)
}

func (s *BetService) GetBets(dealId string) ([]Bet, error) {
	var dealObjectId, err = primitive.ObjectIDFromHex(dealId)
	if err != nil {
		return nil, err
	}

	options := options.Find()
	options.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var result, findErr = s.getBetCollection().Find(s.ctx, bson.D{{
		Key:   "dealID",
		Value: dealObjectId,
	}}, options)
	if findErr != nil {
		return nil, findErr
	}

	var bets = []Bet{}
	var allErr = result.All(s.ctx, &bets)
	if allErr != nil {
		return nil, allErr
	}
	return bets, nil
}
