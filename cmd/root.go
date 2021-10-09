/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ashishbabar/erc20-listener/models"
	"github.com/ashishbabar/erc20-listener/services"
	"github.com/ashishbabar/erc20-listener/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "erc20-listener",
	Short: "This is app listens ERC20 tokens events.",
	Long: `erc20-listener is CLI application implemented using Golang with robust Cobra framework. 
	It takes ERC20 token address and network URL as input and starts listening to 
	events like transfer, approve etc. 
	You can start this listener by passing token address and network URL. For example:

	erc20-listener --contract-address "" --network-url ""
	`,
	Run: start,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.erc20-listener.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".erc20-listener")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func start(cmd *cobra.Command, args []string) {
	logger := util.Zaplogger
	databaseUrl := viper.GetString("database_url")
	ethereumUrl := viper.GetString("network_url")
	contractAddress := viper.GetString("contract_address")

	var eventsToHandle models.Models

	transferEventModel := models.NewTransferEvent(crypto.Keccak256Hash([]byte(viper.GetString("transfer_event_signature"))).Hex(), viper.GetString("transfer_collection_name"))
	approveEventModel := models.NewApproveEvent(crypto.Keccak256Hash([]byte(viper.GetString("approve_event_signature"))).Hex(), viper.GetString("approve_collection_name"))
	eventsToHandle = models.Models{transferEventModel, approveEventModel}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Initializing database client " + databaseUrl)
	databaseClient, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUrl))
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Initializing ethereum client at address " + ethereumUrl)
	chainClient, err := ethclient.Dial(ethereumUrl)
	if err != nil {
		logger.Fatal(err.Error())
	}

	dbClient := util.NewDB(databaseClient)

	listener := services.NewListerner(dbClient, logger, chainClient)

	listener.Start(contractAddress, &eventsToHandle)
}
