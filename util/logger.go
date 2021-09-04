package util

import "go.uber.org/zap"

var Zaplogger *zap.Logger

type Logger struct {
	Logger *zap.Logger
}

func init() {
	// TODO:Read yaml/json from config and build config object
	// var cfg zap.Config
	// if err := json.Unmarshal(rawJSON, &cfg); err != nil {
	// 	panic(err)
	// }
	// logger, err := cfg.Build()
	// if err != nil {
	// 	panic(err)
	// }

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	Zaplogger, _ = config.Build()
}

// This will be used to set a logger with different configuration
func New(logger *zap.Logger) {
	Zaplogger = logger
}
