package vmstorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SnapshotCreatePath defines vmstorage endpoint to create instant snapshot
// https://github.com/VictoriaMetrics/VictoriaMetrics/wiki/Cluster-VictoriaMetrics#url-format
var SnapshotCreatePath = "/snapshot/create"

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

func (c *SnapClient) CreateSnapshot() *SnapResponse {
	resp, err := http.Get(c.ComposeBasePath() + SnapshotCreatePath)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	r := SnapResponse{}
	json.Unmarshal(body, &r)

	return &r
}

func (c SnapClient) ComposeBasePath() string {
	return fmt.Sprintf("%s://%s:%d", c.proto, c.host, c.port)
}
