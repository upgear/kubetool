package kubetool

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func Build(in Input) error {
	name := "docker"
	params := []string{"build", "-t", in.ComputedEnv.ContainerImage, "-f", in.ComputedEnv.DockerFile, in.ComputedEnv.DockerContext}
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
