package util

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type DbClient struct {
	DB *mongo.Client
}

func NewDB(client *mongo.Client) *DbClient {
	return &DbClient{client}
}

func (client *DbClient) TestFunction() {
	fmt.Println(client)
}
