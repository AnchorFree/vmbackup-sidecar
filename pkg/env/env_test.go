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
	errExpected bool
}

var testCases = []envTestCase{
	{
		descr:       "All vars are set; no error",
		host:        "localhost",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		errExpected: false,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", HostVarName),
		host:        "",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", PortVarName),
		host:        "localhost",
		port:        "",
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", BucketVarName),
		host:        "localhost",
		port:        "4242",
		bucketName:  "",
		dataPath:    "/foo/data",
		errExpected: true,
	},
	{
		descr:       fmt.Sprintf("%s is not set; error", DataPathVarName),
		host:        "localhost",
		port:        "4242",
		bucketName:  "foo-bucket",
		dataPath:    "",
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
