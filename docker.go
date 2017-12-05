package kubetool

import (
	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func Build(in CommandInput) error {
	params := []string{"build", "-t", in.Env.ContainerImage, "-f", in.Env.DockerFile, in.Env.DockerContext}

	if _, err := cmd(in.Flags.Verbose, "docker", params...); err != nil {
		return errors.Wrap(err, "unable to build docker image")
	}

	return nil
}

func Push(in CommandInput) error {
	dolog := in.Flags.Verbose

	if in.Flags.Env == DevEnv {
		if dolog {
			log.Info("skipping push because of flag", log.M{"env": in.Flags.Env})
		}
		return nil
	}

	if _, err := cmd(dolog, in.Cloud, "docker", "--", "push", in.ContainerImage); err != nil {
		return errors.Wrap(err, "unable to push docker image")
	}

	return nil
}
