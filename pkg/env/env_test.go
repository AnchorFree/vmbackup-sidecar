package env

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type envTestCase struct {
	descr       string
	host        string
	port        string
	bucketName  string
	dataPath    string
	podName     string
	errExpected bool
}

var testCases = []envTestCase{
	{
		descr:       "All vars are set; no error",
		host:        "localhost",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		podName:     "victoria-metrics-vmstorage-0",
		errExpected: false,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", HostVarName),
		host:        "",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		podName:     "victoria-metrics-vmstorage-0",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", PortVarName),
		host:        "localhost",
		port:        "",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		podName:     "victoria-metrics-vmstorage-0",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", BucketVarName),
		host:        "localhost",
		port:        "4242",
		bucketName:  "",
		dataPath:    "/foo/data",
		podName:     "victoria-metrics-vmstorage-0",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", DataPathVarName),
		host:        "localhost",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "",
		podName:     "victoria-metrics-vmstorage-0",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", PodVarName),
		host:        "localhost",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		podName:     "",
		errExpected: true,
	},
}

func (s envTestCase) getExpectedCfg() *Config {
	var c Config
	if !s.errExpected {
		p, err := strconv.ParseUint(s.port, 10, 16)
		if err != nil {
			panic(err)
		}
		c = Config{
			Host:       s.host,
			Port:       uint16(p),
			BucketName: s.bucketName,
			DataPath:   s.dataPath,
			PodName:    s.podName,
		}
	}
	return &c
}

func TestGetConfig(t *testing.T) {
	for _, tCase := range testCases {
		// Env setup
		os.Setenv(HostVarName, tCase.host)
		os.Setenv(PortVarName, tCase.port)
		os.Setenv(BucketVarName, tCase.bucketName)
		os.Setenv(DataPathVarName, tCase.dataPath)
		os.Setenv(PodVarName, tCase.podName)

		// Actual test
		res, err := GetConfig()
		assert.Equal(t, tCase.getExpectedCfg(), res, tCase.descr)
		if tCase.errExpected {
			assert.NotNil(t, err, tCase.descr)
		} else {
			assert.Nil(t, err, tCase.descr)
		}
	}
}
