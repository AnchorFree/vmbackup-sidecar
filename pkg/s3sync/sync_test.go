package s3sync

import "testing"

type syncCmdTestCase struct {
	bucketName string
	localPath  string
	delete     bool
	follow     bool
	descr      string
	expected   string
}

var testCases = []syncCmdTestCase{
	{
		descr:      "delete + follow",
		bucketName: "vm-backup",
		localPath:  "/foo/bar",
		delete:     true,
		follow:     true,
		expected:   "aws s3 sync /foo/bar s3://vm-backup --delete --follow-symlinks",
	},
	{
		descr:      "no delete, no follow",
		bucketName: "vm-backup",
		localPath:  "/foo/bar",
		delete:     false,
		follow:     false,
		expected:   "aws s3 sync /foo/bar s3://vm-backup",
	},
}

func TestComposeCmd(t *testing.T) {
	for _, tCase := range testCases {
		res := New(
			tCase.bucketName,
			tCase.localPath,
			tCase.delete,
			tCase.follow,
		).ComposeCmd()
		if res != tCase.expected {
			t.Fatalf("FAIL: %s\nExpected: '%s'\nActual: '%s'", tCase.descr, tCase.expected, res)
		}
		t.Logf("PASS: %s", tCase.descr)
	}
}
