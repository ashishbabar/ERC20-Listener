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

type TransferEvent struct {
	eventSignature string
	collectionName string
}

func NewTransferEvent(_eventSignature string, _collectionName string) *TransferEvent {
	return &TransferEvent{
		eventSignature: _eventSignature,
		collectionName: _collectionName,
	}
}
func (event *TransferEvent) GetConfig() *Config {
	return &Config{event.eventSignature, event.collectionName}
}
func (event *TransferEvent) Store(collection *mongo.Collection, eventLog *types.Log, contractAbi *abi.ABI, logger *zap.Logger) bool {
	logger.Info("Received event : Transfer")

	if eventLog.Removed {
		return false
	}
	var transferEvent token.SimpleTokenTransfer

	err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", eventLog.Data)
	if err != nil {
		logger.Fatal(err.Error())
	}

	transferEvent.From = common.HexToAddress(eventLog.Topics[1].Hex())
	transferEvent.To = common.HexToAddress(eventLog.Topics[2].Hex())

	logger.Info("From: " + transferEvent.From.Hex())
	logger.Info("To: " + transferEvent.To.Hex())
	logger.Info("Tokens: " + transferEvent.Value.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "from", Value: transferEvent.From.Hex()}, {Key: "to", Value: transferEvent.To.Hex()}, {Key: "amount", Value: transferEvent.Value.String()}, {Key: "created_at", Value: time.Now()}})
	if err != nil {
		logger.Fatal(err.Error())
		return false
	}
	logger.Info("Inserted transfer event in database" + fmt.Sprintf("%v", res.InsertedID))
	return true
}
