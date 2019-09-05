package vmstorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// https://github.com/VictoriaMetrics/VictoriaMetrics/wiki/Cluster-VictoriaMetrics#url-format

	// SnapshotCreatePath defines vmstorage endpoint to create instant snapshot
	SnapshotCreatePath = "/snapshot/create"

	// SnapshotDeleteAll defines vmstorage endpoint to delete all snapshots
	SnapshotDeleteAll = "/snapshot/delete_all"
)

type SnapClient struct {
	proto string // default: "http"
	host  string
	port  uint16
}

type SnapResponse struct {
	Status   string `json:"status"`
	SnapName string `json:"snapshot"`
}

// New creates an instance of SnapClient
func New(host string, port uint16, proto string) *SnapClient {
	if proto == "" {
		proto = "http"
	}
	return &SnapClient{host: host, port: port, proto: proto}
}

func (c SnapClient) CreateSnapshot() (*SnapResponse, error) {
	return c.makeApiRequest(SnapshotCreatePath)
}

func (c SnapClient) DeleteAllSnaps() (*SnapResponse, error) {
	return c.makeApiRequest(SnapshotDeleteAll)
}

func (c SnapClient) ComposeBasePath() string {
	return fmt.Sprintf("%s://%s:%d", c.proto, c.host, c.port)
}

func (c SnapClient) makeApiRequest(path string) (*SnapResponse, error) {
	var r SnapResponse

	resp, err := http.Get(c.ComposeBasePath() + path)
	if err != nil {
		return &r, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &r, err
	}

	json.Unmarshal(body, &r)
	return &r, nil
}
