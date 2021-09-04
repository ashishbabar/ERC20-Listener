module github.com/ashishbabar/erc20-listener

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.8
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	go.mongodb.org/mongo-driver v1.7.1
	go.uber.org/zap v1.17.0
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/gorm v1.21.13
)

replace github.com/ashishbabar/erc20-listener/contracts => ./contracts
