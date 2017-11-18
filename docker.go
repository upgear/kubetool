package kubetool

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
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

func Push(in Input) error {
	dolog := in.Flags.Verbose

	if in.Flags.Local {
		if dolog {
			log.Info("skipping push because of flag", log.M{"flag": "local"})
		}
		return nil
	}

	if _, err := cmd(dolog, in.ComputedEnv.Cloud, "docker", "--", "push", in.ComputedEnv.ContainerImage); err != nil {
		return errors.Wrap(err, "unable to push docker image")
	}

	return nil
}
