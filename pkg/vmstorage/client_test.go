package vmstorage

import "testing"

type clientTestCase struct {
	descr    string
	proto    string
	host     string
	port     uint16
	expected *SnapClient
	basePath string
}

var testCases = []clientTestCase{
	{
		descr:    "default proto value",
		proto:    "",
		host:     "foo",
		port:     8899,
		expected: &SnapClient{host: "foo", port: 8899, proto: "http"},
		basePath: "http://foo:8899",
	},
	{
		descr:    "https proto value",
		proto:    "https",
		host:     "bar",
		port:     8080,
		expected: &SnapClient{host: "bar", port: 8080, proto: "https"},
		basePath: "https://bar:8080",
	},
}

func TestNew(t *testing.T) {
	for _, tCase := range testCases {
		res := New(tCase.host, tCase.port, tCase.proto)
		if *res != *tCase.expected {
			t.Fatalf("FAIL: %s\nExpected: %+v\nActual: %+v", tCase.descr, tCase.expected, res)
		}
		t.Logf("PASS: %s", tCase.descr)
	}
}

func TestComposeBasePath(t *testing.T) {
	for _, tCase := range testCases {
		c := New(tCase.host, tCase.port, tCase.proto)
		res := c.ComposeBasePath()
		if res != tCase.basePath {
			t.Fatalf("FAIL: %s\nExpected: %+v\nActual: %+v", tCase.descr, tCase.basePath, res)
		}
		t.Logf("PASS: %s", tCase.descr)
	}
}
