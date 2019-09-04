package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/cfg"
	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/handlers"
)

// version and semver get overwritten by build with
// go build -i -v -ldflags="-X main.version=$(git describe --always --long) -X main.semver=v$(git semver get)"
var (
	version = "undefined"
	builddt = "undefined"
	semver  = "undefined"
	branch  = "undefined"
	logger  = cfg.Cfg.Logger
	fastlog = cfg.Cfg.FastLogger
)

func main() {
	listen := ":" + strconv.Itoa(cfg.Cfg.Port)
	logger.Infow(
		"starting "+os.Args[0],
		"version", version,
		"buildtime", builddt,
		"semver", semver,
		"branch", branch,
		"listen", listen,
	)
	// TODO: implement /metrics
	// http.HandleFunc("/metrics", metricsHandler)

	http.HandleFunc("/health", handlers.HealthcheckHandler)
	http.HandleFunc("/backup/create", handlers.BackupHandler)

	log.Fatal(http.ListenAndServe(listen, nil))
}
