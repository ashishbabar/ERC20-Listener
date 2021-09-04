package models

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type TransferEvent struct {
	from   string   `bson:"from,omitempty"`
	to     string   `bson:"to,omitempty"`
	amount *big.Int `bson:"amount,omitempty"`
}

func NewTransferEvent(_from string, _to string, _amount *big.Int) *TransferEvent {
	return &TransferEvent{from: _from, to: _to, amount: _amount}
}
func (event *TransferEvent) Store(collection *mongo.Collection, logger *zap.Logger) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "from", Value: event.from}, {Key: "to", Value: event.to}, {Key: "amount", Value: event.amount.String()}})
	if err != nil {
		logger.Fatal(err.Error())
		return false
	}
	id := res.InsertedID
	// if str, ok := id.(string); ok {
	// 	logger.Info("Inserted transfer event in database" + str)
	// } else {
	// 	logger.Error("Unable to parse ID returned after insertion")
	// }
	logger.Info("Inserted transfer event in database" + fmt.Sprintf("%v", id))
	return true
}
