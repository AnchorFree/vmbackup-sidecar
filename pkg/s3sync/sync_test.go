package s3sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(t, tCase.expected, res, tCase.descr)
	}
}
