package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/cfg"
	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/handlers"
	"go.uber.org/zap"
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

// fwdToZapWriter allows us to use the zap.Logger as our http.Server ErrorLog
// see https://stackoverflow.com/questions/52294334/net-http-set-custom-logger
type fwdToZapWriter struct {
	logger *zap.Logger
}

func (fw *fwdToZapWriter) Write(p []byte) (n int, err error) {
	fw.logger.Error(string(p))
	return len(p), nil
}

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

	discardBackupHandler := handlers.DiscardConcRequests(
		handlers.BackupHandler,
		"Backup ongoing, discarding request",
		http.StatusConflict,
	)
	http.HandleFunc("/backup/create", discardBackupHandler)

	srv := &http.Server{
		Addr: listen,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     log.New(&fwdToZapWriter{fastlog}, "", 0),
	}
	log.Fatal(srv.ListenAndServe())
}
