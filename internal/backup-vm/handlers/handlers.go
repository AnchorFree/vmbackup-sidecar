package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/AnchorFree/vmbackup-sidecar/pkg/env"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/s3sync"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/vmstorage"
)

func BackupHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintln(w, "Creating snapshot")
		client := vmstorage.New(conf.Host, conf.Port, "http")
		resp := client.CreateSnapshot()
		if resp.Status != "ok" {
			fmt.Printf("Error: %s: /snapshot/create status not 'ok'\n", resp.Status)
			os.Exit(1)
		}
		fmt.Fprintf(w, "Snapshot '%s' created\n", resp.SnapName)

		// Sync snapshot with S3
		snapPath := path.Join(conf.DataPath, "snapshots", resp.SnapName)
		bucketPath := path.Join(conf.BucketName, conf.Host)
		delete := true
		follow := true

		fmt.Fprintf(w, "Sync snapshot %s into %s\n", resp.SnapName, bucketPath)
		syncer := s3sync.New(
			bucketPath,
			snapPath,
			delete,
			follow,
		)
		out, err := syncer.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(out))
		fmt.Fprintln(w, "Sync completed")

	default:
		fmt.Fprintln(w, "Error: only PUT method is supported")
	}
}

// func metricsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, ";;; call to /metrics")
// }

// HealthcheckHandler /health
// just returns 200 '{ "ok": true }'
func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, "{ \"ok\": true }"); err != nil {
		fmt.Printf("Error in response writing: %#v", err)
	}
}
