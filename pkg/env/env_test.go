package env

import (
	"os"
	"strconv"
	"testing"
)

type envTestCase struct {
	descr       string
	host        string
	port        uint16
	bucketName  string
	dataPath    string
	expectedErr error
}

var testCases = []envTestCase{
	{
		descr:       "All env vars are set",
		host:        "localhost",
		port:        4242,
		bucketName:  "foo-bucket",
		dataPath:    "/foo/data",
		expectedErr: nil,
	},
}

func (s envTestCase) getExpectedCfg() *Config {
	return &Config{
		Host:       s.host,
		Port:       s.port,
		BucketName: s.bucketName,
		DataPath:   s.dataPath,
	}
}

func TestGetConfig(t *testing.T) {
	for _, tCase := range testCases {
		// Env setup
		os.Setenv(HostVarName, tCase.host)
		os.Setenv(PortVarName, strconv.Itoa(int(tCase.port)))
		os.Setenv(BucketVarName, tCase.bucketName)
		os.Setenv(DataPathVarName, tCase.dataPath)

		// Actual test
		res, err := GetConfig()
		expectedCfg := tCase.getExpectedCfg()
		if *res != *expectedCfg {
			t.Fatalf("FAIL: %s\nExpected: %+v\nActual: %+v", tCase.descr, expectedCfg, res)
		}
		if err != tCase.expectedErr {
			t.Fatalf("FAIL: %s\nExpected: '%s'\nActual: '%s'", tCase.descr, tCase.expectedErr, err)
		}
		t.Logf("PASS: %s", tCase.descr)
	}
}
