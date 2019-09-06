package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/AnchorFree/vmbackup-sidecar/internal/backup-vm/cfg"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/env"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/s3sync"
	"github.com/AnchorFree/vmbackup-sidecar/pkg/vmstorage"
)

var (
	log             = cfg.Cfg.Logger
	envConf, envErr = env.GetConfig()
)

func init() {
	if envErr != nil {
		log.Errorw("error parsing envConfig from env", "error", envErr)
		os.Exit(1)
	}
	log.Infow(
		"configuration from env",
		env.HostVarName, envConf.Host,
		env.PortVarName, envConf.Port,
		env.BucketVarName, envConf.BucketName,
		env.DataPathVarName, envConf.DataPath,
		env.PodVarName, envConf.PodName,
	)
}

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

	// Create snapshot
	fmt.Fprintln(w, "Creating snapshot")
	client := vmstorage.New(envConf.Host, envConf.Port, "http")
	createResp, err := client.CreateSnapshot()
	if err != nil {
		errMsg := "error creating vmstorage snapshot"
		log.Errorw(errMsg, "error", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	if createResp.Status != "ok" {
		errMsg := fmt.Sprintf(
			"vmstorage %s response status not 'ok'",
			vmstorage.SnapshotCreatePath,
		)
		log.Errorw(errMsg, "status", createResp.Status)
		http.Error(w, "failed to create snapshot", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Snapshot '%s' created\n", createResp.SnapName)

	// Sync snapshot with S3
	snapPath := path.Join(envConf.DataPath, "snapshots", createResp.SnapName)
	bucketPath := path.Join(envConf.BucketName, envConf.PodName)
	delete := true
	follow := true

	fmt.Fprintf(w, "Sync snapshot %s into %s\n", createResp.SnapName, bucketPath)
	syncer := s3sync.New(bucketPath, snapPath, delete, follow)
	out, err := syncer.Run()
	if err != nil {
		log.Errorw("error syncing snapshot with s3", "error", err)
		http.Error(w, "failed to sync snapshot with s3", http.StatusInternalServerError)
		return
	}
	log.Info(string(out))
	fmt.Fprintln(w, "Sync completed")

	// Remove all snapshots
	fmt.Fprintln(w, "Removing all snapshots")
	delAllResp, err := client.DeleteAllSnaps()
	if err != nil {
		errMsg := "error removing vmstorage snapshots"
		log.Errorw(errMsg, "error", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	if delAllResp.Status != "ok" {
		errMsg := fmt.Sprintf(
			"vmstorage %s response status not 'ok'",
			vmstorage.SnapshotDeleteAll,
		)
		log.Errorw(errMsg, "status", delAllResp.Status)
		http.Error(w, "failed to remove snapshots", http.StatusInternalServerError)
		return
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

// DiscardConcRequests is HTTP handling middleware that ensures no more than
// single request is passed concurrently to the given handler f. Other requests
// are discarded.
func DiscardConcRequests(f http.HandlerFunc, errMsg string, httpStatusCode int) http.HandlerFunc {
	// XXX: very important that channel is buffered
	sema := make(chan struct{}, 1)

	return func(w http.ResponseWriter, req *http.Request) {
		select {
		case sema <- struct{}{}:
			defer func() { <-sema }()
			f(w, req)
		default:
			http.Error(w, errMsg, httpStatusCode)
		}
	}
}
