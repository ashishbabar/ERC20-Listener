module github.com/ashishbabar/erc20-listener

go 1.16

replace github.com/ashishbabar/erc20-listener/cmd => ./cmd

require (
	github.com/ashishbabar/erc20-listener/cmd v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.17.0
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/gorm v1.21.13 // indirect
)

replace github.com/ashishbabar/erc20-listener/util => ./util
