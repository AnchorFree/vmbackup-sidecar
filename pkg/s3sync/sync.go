package s3sync

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

type SyncCmd struct {
	BucketURI      string
	LocalPath      string
	Delete         bool
	FollowSymlinks bool
}

func New(bucketName, localPath string, delete, follow bool) *SyncCmd {
	return &SyncCmd{
		BucketURI:      "s3://" + bucketName,
		LocalPath:      localPath,
		Delete:         delete,
		FollowSymlinks: follow,
	}
}

func (sc SyncCmd) ComposeCmd() string {
	base := "aws s3 sync"
	src := sc.LocalPath
	dst := sc.BucketURI

	delete := ""
	if sc.Delete {
		delete = "--delete"
	}

	follow := ""
	if sc.FollowSymlinks {
		follow = "--follow-symlinks"
	}
	cmd := fmt.Sprintf("%s %s %s %s %s", base, src, dst, delete, follow)
	return strings.TrimSpace(cmd)
}

// Run executes composed `aws s3 sync` command.
func (sc SyncCmd) Run() ([]byte, error) {
	return runCmd(sc.ComposeCmd())
}

// KeepEmptyDirs fills empty directories with .keep file.
// This is useful for syncing empty directories in S3.
// S3 is an object storage which does not have a concept
// of the empty directory.
// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
func KeepEmptyDirs(path string) ([]byte, error) {
	cmd := fmt.Sprintf("find %s -follow -type d -empty -exec touch {}/.keep \\;", path)
	return runCmd(cmd)
}

func runCmd(cmd string) ([]byte, error) {
	var out []byte

	tokens, err := shlex.Split(cmd)
	if err != nil {
		return out, err
	}

	proc := exec.Command(tokens[0], tokens[1:]...)
	out, err = proc.Output()
	if err != nil {
		return out, err
	}

	return out, err
}
