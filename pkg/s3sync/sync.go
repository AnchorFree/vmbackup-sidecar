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
	Profile        string
	Delete         bool
	FollowSymlinks bool
}

func New(bucketName, localPath, profile string, delete, follow bool) *SyncCmd {
	return &SyncCmd{
		BucketURI:      "s3://" + bucketName,
		LocalPath:      localPath,
		Profile:        profile,
		Delete:         delete,
		FollowSymlinks: follow,
	}
}

func (sc SyncCmd) ComposeCmd() string {
	base := "aws s3 sync"
	src := sc.LocalPath
	dst := sc.BucketURI
	profile := "--profile " + sc.Profile

	delete := ""
	if sc.Delete {
		delete = "--delete"
	}

	follow := ""
	if sc.FollowSymlinks {
		follow = "--follow-symlinks"
	}
	cmd := fmt.Sprintf("%s %s %s %s %s %s", base, profile, src, dst, delete, follow)
	return strings.TrimSpace(cmd)
}

func (sc SyncCmd) Run() ([]byte, error) {
	var out []byte

	cmd := sc.ComposeCmd()
	fmt.Printf(";;; cmd: '%s'\n", cmd)

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
