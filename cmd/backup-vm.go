package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/handlers"
)

func main() {
	http.HandleFunc("/backup", handlers.BackupHandler)
	listenPort := "8488"
	fmt.Printf(";;; Listening port %s\n", listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, nil))

	// TODO: implement /metrics
	// http.HandleFunc("/metrics", metricsHandler)

	// TODO: implement /health
	// http.HandleFunc("/health", healthHandler)
}
