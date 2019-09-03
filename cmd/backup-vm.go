package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/handlers"
)

func main() {
	// TODO: implement /metrics
	// http.HandleFunc("/metrics", metricsHandler)

	http.HandleFunc("/health", handlers.HealthcheckHandler)
	http.HandleFunc("/backup/create", handlers.BackupHandler)

	listenPort := "8488"
	fmt.Printf(";;; Listening port %s\n", listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, nil))
}
