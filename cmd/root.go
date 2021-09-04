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

	"github.com/ashishbabar/erc20-listener/services"
	"github.com/ashishbabar/erc20-listener/util"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "erc20-listener",
	Short: "This is app listens ERC20 tokens events.",
	Long: `erc20-listener is CLI application implemented using Golang with robust Cobra framework. 
	It takes ERC20 token address and network URL as input and starts listening to 
	events like transfer, approve etc. 
	You can start this listener by passing token address and network URL. For example:

	erc20-listener --contract-address "" --network-url ""
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: start,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.erc20-listener.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".erc20-listener" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".erc20-listener")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func start(cmd *cobra.Command, args []string) {
	logger := util.Zaplogger
	databaseUrl := viper.GetString("database_url")
	ethereumUrl := viper.GetString("network_url")
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

	listener.Start(viper.GetString("contract_address"))
}
