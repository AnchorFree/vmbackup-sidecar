package env

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// vmstorage host and port
	Host string
	Port uint16

	// S3 bucket name for syncing snapshot
	BucketName string

	// Correspondes to --storageDataPath flag in VictoriaMetrics setup
	DataPath string

	// awscli profile with access to provided S3 bucket
	Profile string
}

func GetConfig() (*Config, error) {
	var s Config

	hostVarName := "VMSTORAGE_HOST"
	portVarName := "VMSTORAGE_PORT"
	bucketVarName := "VM_SNAPSHOT_BUCKET"
	dataPathVarName := "VM_STORAGE_DATA_PATH"
	profileVarName := "VM_AWS_PROFILE"

	host := os.Getenv(hostVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", hostVarName)
	}

	port := os.Getenv(portVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", portVarName)
	}
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return &s, err
	}

	bucket := os.Getenv(bucketVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", bucketVarName)
	}

	dataPath := os.Getenv(dataPathVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", dataPathVarName)
	}

	profile := os.Getenv(profileVarName)
	if host == "" {
		return &s, fmt.Errorf("%s is not set", profileVarName)
	}

	s = Config{
		Host:       host,
		Port:       uint16(p),
		BucketName: bucket,
		DataPath:   dataPath,
		Profile:    profile,
	}

	return &s, nil
}
