package kubetool

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Build(in Input) error {
	tag, err := templateString(in.Env.TagTemplate, in)
	if err != nil {
		return errors.Wrap(err, "unable to template tag")
	}

	file := filepath.Join(in.Env.DockerfilePath, fmt.Sprintf("%s.Dockerfile", in.Args.Name))

	name := "docker"
	params := []string{"build", "-t", tag, "-f", file, in.Env.DockerContext}
	cmd := exec.Command(name, params...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if in.Flags.Verbose {
		cmd.Stdout = os.Stdout
		logCmd(name, params...)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "unable to execute command: %s %s: %s", name, strings.Join(params, " "), stderr.String())
	}

	return nil
}
