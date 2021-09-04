package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	token "github.com/ashishbabar/erc20-listener/contracts"
	"github.com/ashishbabar/erc20-listener/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var transferFuncSig = []byte("Transfer(address,address,uint256)")
var approvalFunSig = []byte("Approval(address,address,uint256)")
var transferFuncSigHash = crypto.Keccak256Hash(transferFuncSig)
var approvalFunSigHash = crypto.Keccak256Hash(approvalFunSig)

type Listener struct {
	dbClient  *util.DbClient
	zapLogger *zap.Logger
	ethClient *ethclient.Client
}

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func NewListerner(client *util.DbClient, logger *zap.Logger, chainClient *ethclient.Client) *Listener {
	return &Listener{dbClient: client, zapLogger: logger, ethClient: chainClient}
}

func (listener *Listener) Start(contractAddress string) {
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
			handleEventLog(&vLog, logger, &contractAbi)
		}
	}

}

func handleEventLog(vLog *types.Log, logger *zap.Logger, contractAbi *abi.ABI) {

	logger.Info("Log Block Number: " + strconv.FormatUint(uint64(vLog.BlockNumber), 10))
	logger.Info("Log Index: " + strconv.FormatUint(uint64(vLog.Index), 10))

	switch vLog.Topics[0].Hex() {
	case transferFuncSigHash.Hex():
		//
		logger.Info("Received event : Transfer")

		var transferEvent token.SimpleTokenTransfer

		err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
		if err != nil {
			logger.Fatal(err.Error())
		}

		transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
		transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

		logger.Info("From: " + transferEvent.From.Hex())
		logger.Info("To: " + transferEvent.To.Hex())
		logger.Info("Tokens: " + transferEvent.Value.String())
	case approvalFunSigHash.Hex():
		//
	}
}
