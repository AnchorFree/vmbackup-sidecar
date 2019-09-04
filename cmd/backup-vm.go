package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	logger.Infow(
		"starting "+os.Args[0],
		"version", version,
		"buildtime", builddt,
		"semver", semver,
		"branch", branch,
	)
	// TODO: implement /metrics
	// http.HandleFunc("/metrics", metricsHandler)

	http.HandleFunc("/health", handlers.HealthcheckHandler)
	http.HandleFunc("/backup/create", handlers.BackupHandler)

	listenPort := "8488"
	fmt.Printf(";;; Listening port %s\n", listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, nil))
}
