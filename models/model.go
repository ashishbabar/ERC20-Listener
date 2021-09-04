package models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Model interface {
	Store(*mongo.Collection, *zap.Logger) bool
}
