package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	token "github.com/ashishbabar/erc20-listener/contracts"
	"github.com/ashishbabar/erc20-listener/models"
	"github.com/ashishbabar/erc20-listener/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Listener struct {
	dbClient  *util.DbClient
	zapLogger *zap.Logger
	ethClient *ethclient.Client
}

func NewListerner(client *util.DbClient, logger *zap.Logger, chainClient *ethclient.Client) *Listener {
	return &Listener{dbClient: client, zapLogger: logger, ethClient: chainClient}
}

func (listener *Listener) Start(contractAddress string, eventsToHandle *models.Models) {
	logger := listener.zapLogger
	logger.Info("Now starting listening to contract " + contractAddress)
	contractHexAddress := common.HexToAddress(contractAddress)
	logger.Info("Preparing filter query")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractHexAddress},
	}
	logs := make(chan types.Log)
	logger.Info("Started filtering Block logs for events")
	sub, err := listener.ethClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Logs %v", logs)

	logger.Info("Retrieving contract ABI from GO binding")
	contractAbi, err := abi.JSON(strings.NewReader(string(token.SimpleTokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Intializing contract function signatures")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			handleEventLog(&vLog, logger, &contractAbi, listener.dbClient, eventsToHandle)
		}
	}

}

func handleEventLog(vLog *types.Log, logger *zap.Logger, contractAbi *abi.ABI, dbClient *util.DbClient, eventsToHandle *models.Models) {

	logger.Info("Log Block Number: " + strconv.FormatUint(uint64(vLog.BlockNumber), 10))
	logger.Info("Log Index: " + strconv.FormatUint(uint64(vLog.Index), 10))

	databaseObj := dbClient.DB.Database(viper.GetString("database_name"))
	for _, p := range *eventsToHandle {
		config := p.GetConfig()
		if config.EventSignature == vLog.Topics[0].Hex() {
			collection := databaseObj.Collection(config.CollectionName)
			p.Store(collection, vLog, contractAbi, logger)
		}
	}
}
