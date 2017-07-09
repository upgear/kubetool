package kubetool

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func CheckRepo(in *Input) error {
	dolog := in.Flags.Verbose

	// Make sure there all unstaged changes are committed
	lsFiles, err := cmd(dolog, "git", "ls-files",
		"--exclude-standard",
		"--others",
		"--modified",
		"--deleted")
	if err != nil {
		return err
	}
	if lsFiles != "" {
		return errors.New("uncommitted unstaged git changes")
	}

	// Make sure there all staged changes are committed
	stagedDiff, err := cmd(dolog, "git", "--no-pager", "diff", "--cached")
	if err != nil {
		return err
	}
	if stagedDiff != "" {
		return errors.New("uncommitted staged git changes")
	}

	// Make sure we are on the master branch
	ref, err := cmd(dolog, "git", "symbolic-ref", "HEAD")
	if err != nil {
		return err
	}
	if ref != "refs/heads/master" {
		return errors.New("not on master branch")
	}

	// Make sure local master is in sync with remote
	return repoIsSynced(in)
}

var refrgx = regexp.MustCompile(`(\S+)\s+(\S+)`)

func repoIsSynced(in *Input) error {
	dolog := in.Flags.Verbose

	local, err := cmd(dolog, "git", "show-ref", "--heads")
	if err != nil {
		return err
	}

	remote, err := cmd(dolog, "git", "ls-remote", "--heads", "origin")
	if err != nil {
		return err
	}

	refmap := func(refs []string) map[string]string {
		m := make(map[string]string)
		for _, r := range refs {
			matches := refrgx.FindStringSubmatch(r)
			if len(matches) != 3 {
				continue
			}
			m[matches[2]] = matches[1]
		}
		return m
	}

	localM := refmap(strings.Split(local, "\n"))
	remoteM := refmap(strings.Split(remote, "\n"))

	const master = "refs/heads/master"
	if _, ok := localM[master]; !ok {
		return errors.New("local master branch not found")
	}
	if _, ok := remoteM[master]; !ok {
		return errors.New("remote master branch not found")
	}
	if localM[master] != remoteM[master] {
		return fmt.Errorf("remote master ref (%s) not in sync with local (%s)", localM[master], remoteM[master])
	}

	in.Repo.CommitHash = localM[master]

	return nil
}
