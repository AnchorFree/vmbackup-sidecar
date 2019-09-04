package cfg

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type config struct {
	Logger     *zap.SugaredLogger
	FastLogger *zap.Logger
	Port       int
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

	// Args from cmd
	port := flag.Int("port", 8488, "Port to listen")
	help := flag.Bool("help", false, "Show usage")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(2)
	}

	Cfg.Port = *port
}
