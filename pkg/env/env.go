package env

import (
	"fmt"
	"os"
	"strconv"
)

// Env vars to read
var (
	HostVarName     = "VMSTORAGE_HOST"
	PortVarName     = "VMSTORAGE_PORT"
	BucketVarName   = "VM_SNAPSHOT_BUCKET"
	DataPathVarName = "VM_STORAGE_DATA_PATH"
)

type Config struct {
	// vmstorage host and port
	Host string
	Port uint16

	// S3 bucket name for syncing snapshot
	BucketName string

	// Correspondes to --storageDataPath flag in VictoriaMetrics setup
	DataPath string
}

func GetConfig() (*Config, error) {
	var s Config

	host := os.Getenv(HostVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", HostVarName)
	}

	port := os.Getenv(PortVarName)
	if port == "" {
		return &s, fmt.Errorf("%s is not set", PortVarName)
	}
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return &s, err
	}

	bucket := os.Getenv(BucketVarName)
	if bucket == "" {
		return &s, fmt.Errorf("%s is not set", BucketVarName)
	}

	dataPath := os.Getenv(DataPathVarName)
	if dataPath == "" {
		return &s, fmt.Errorf("%s is not set", DataPathVarName)
	}

	s = Config{
		Host:       host,
		Port:       uint16(p),
		BucketName: bucket,
		DataPath:   dataPath,
	}

	return &s, nil
}
