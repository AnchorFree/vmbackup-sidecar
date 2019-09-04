package cfg

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

type config struct {
	Logger     *zap.SugaredLogger
	FastLogger *zap.Logger
	// LogLevel   string
	// Listen     string
	// Port       int
}

var Cfg config

func init() {
	// Setup logger(s)
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	Cfg.Logger = sugar
	Cfg.FastLogger = logger
}
