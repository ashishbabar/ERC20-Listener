package models

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Config struct {
	EventSignature string
	CollectionName string
}
type Model interface {
	GetConfig() *Config
	Store(*mongo.Collection, *types.Log, *abi.ABI, *zap.Logger) bool
}

type Models []Model
