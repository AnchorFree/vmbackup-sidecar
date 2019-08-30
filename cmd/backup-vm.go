package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/AnchorFree/vmbackup-sidecar/pkg/env"

	"github.com/AnchorFree/vmbackup-sidecar/pkg/s3sync"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/vmstorage"
)

func main() {
	http.HandleFunc("/backup", backupHandler)
	listenPort := "8488"
	fmt.Printf(";;; Listening port %s\n", listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, nil))

	// TODO: implement /metrics
	// http.HandleFunc("/metrics", metricsHandler)

	// TODO: implement /health
	// http.HandleFunc("/health", healthHandler)
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ";;; call to /backup\n")

	switch r.Method {
	case "PUT":
		// Read ENV vars
		conf, err := env.GetConfig()
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		// Create snapshot
		fmt.Fprintf(w, "Creating snapshot")
		client := vmstorage.New(conf.Host, conf.Port, "http")
		resp := client.CreateSnapshot()
		if resp.Status != "ok" {
			fmt.Printf("Error: %s: /snapshot/create status not 'ok'\n", resp.Status)
			os.Exit(1)
		}
		fmt.Fprintf(w, "Snapshot '%s' created", resp.SnapName)

		// Sync snapshot with S3
		snapPath := path.Join(conf.DataPath, "snapshots", resp.SnapName)
		bucketPath := path.Join(conf.BucketName, conf.Host)
		delete := true
		follow := true

		fmt.Fprintf(w, "Sync snapshot %s into %s", resp.SnapName, bucketPath)
		syncer := s3sync.New(
			bucketPath,
			snapPath,
			conf.Profile,
			delete,
			follow,
		)
		out, err := syncer.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(out))

	default:
		fmt.Fprintf(w, "Error: only PUT method is supported")
	}
}

// func metricsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, ";;; call to /metrics")
// }

// func healthHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, ";;; call to /health")
// }
