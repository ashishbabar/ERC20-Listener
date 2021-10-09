package models

import (
	"context"
	"fmt"
	"time"

	token "github.com/ashishbabar/erc20-listener/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ApproveEvent struct {
	eventSignature string
	collectionName string
}

func NewApproveEvent(_eventSignature string, _collectionName string) *ApproveEvent {
	return &ApproveEvent{
		eventSignature: _eventSignature,
		collectionName: _collectionName,
	}
}

func (event *ApproveEvent) GetConfig() *Config {
	return &Config{event.eventSignature, event.collectionName}
}

func (event *ApproveEvent) Store(collection *mongo.Collection, eventLog *types.Log, contractAbi *abi.ABI, logger *zap.Logger) bool {
	logger.Info("Received event : Approve")

	if eventLog.Removed {
		return false
	}
	var approveEvent token.SimpleTokenApproval

	err := contractAbi.UnpackIntoInterface(&approveEvent, "Approve", eventLog.Data)
	if err != nil {
		logger.Fatal(err.Error())
	}

	approveEvent.Owner = common.HexToAddress(eventLog.Topics[1].Hex())
	approveEvent.Spender = common.HexToAddress(eventLog.Topics[2].Hex())

	logger.Info("Owner: " + approveEvent.Owner.Hex())
	logger.Info("Spender: " + approveEvent.Spender.Hex())
	logger.Info("Tokens: " + approveEvent.Value.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "owner", Value: approveEvent.Owner.Hex()}, {Key: "spender", Value: approveEvent.Spender.Hex()}, {Key: "amount", Value: approveEvent.Value.String()}, {Key: "created_at", Value: time.Now()}})
	if err != nil {
		logger.Fatal(err.Error())
		return false
	}
	logger.Info("Inserted transfer event in database" + fmt.Sprintf("%v", res.InsertedID))
	return true
}
