package handlers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/cfg"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/env"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/s3sync"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/vmstorage"
)

var log = cfg.Cfg.Logger

func BackupHandler(w http.ResponseWriter, r *http.Request) {
	pattern := "/backup/create"
	log.Infof("Call to %s", pattern)

	if r.Method != "GET" {
		w.Header().Set("Content-Type", "application/json")
		errMsg := fmt.Sprintf("Incorrect HTTP method for uri [%s] and method [%s], allowed: [GET]", pattern, r.Method)
		errFull := fmt.Sprintf("{ \"error\": \"%s\", \"status\": 405 }", errMsg)
		log.Error(errMsg)
		http.Error(w, errFull, http.StatusMethodNotAllowed)
		return
	}

	// Read ENV vars
	conf, err := env.GetConfig()
	if err != nil {
		errMsg := "error parsing config from env"
		log.Errorw(errMsg, "error", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Create snapshot
	fmt.Fprintln(w, "Creating snapshot")
	client := vmstorage.New(conf.Host, conf.Port, "http")
	resp := client.CreateSnapshot()
	if resp.Status != "ok" {
		errMsg := fmt.Sprintf(
			"vmstorage %s response status not 'ok'",
			vmstorage.SnapshotCreatePath,
		)
		log.Errorw(errMsg, "status", resp.Status)
		http.Error(w, "failed to create snapshot", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Snapshot '%s' created\n", resp.SnapName)

	// Sync snapshot with S3
	snapPath := path.Join(conf.DataPath, "snapshots", resp.SnapName)
	bucketPath := path.Join(conf.BucketName, conf.PodName)
	delete := true
	follow := true

	fmt.Fprintf(w, "Sync snapshot %s into %s\n", resp.SnapName, bucketPath)
	syncer := s3sync.New(bucketPath, snapPath, delete, follow)
	out, err := syncer.Run()
	if err != nil {
		log.Errorw("error syncing snapshot with s3", "error", err)
		http.Error(w, "failed to sync snapshot with s3", http.StatusInternalServerError)
		return
	}
	log.Info(string(out))
	fmt.Fprintln(w, "Sync completed")
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
