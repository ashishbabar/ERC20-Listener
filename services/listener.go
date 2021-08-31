package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	token "github.com/ashishbabar/erc20-listener/contracts"
	"github.com/ashishbabar/erc20-listener/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

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
	// simpleToken, err := token.NewSimpleToken(contractHexAddress, listener.ethClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	logger.Info("Preparing filter query")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6383820),
		ToBlock:   big.NewInt(6383840),
		Addresses: []common.Address{
			contractHexAddress,
		},
	}
	logger.Info("Started filtering Block logs for events")
	logs, err := listener.ethClient.FilterLogs(context.Background(), query)
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
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	sugar := logger.Sugar()
	sugar.Infow("logs", logs)
	logger.Info("Reading logs from response")
	for _, vLog := range logs {

		sugar.Infow("Log Block Number: ", vLog.BlockNumber)
		sugar.Infow("Log Index: %d\n", vLog.Index)
		// fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		// fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			//
			sugar.Infow("Log Name: Transfer\n")
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			sugar.Infow("From: %s\n", transferEvent.From.Hex())
			sugar.Infow("To: %s\n", transferEvent.To.Hex())
			sugar.Infow("Tokens: %s\n", transferEvent.Tokens.String())
		case logApprovalSigHash.Hex():
			//
		}
	}
}
