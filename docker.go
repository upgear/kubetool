package kubetool

import (
	"os"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func Build(in CommandInput) error {
	for i := range in.Env.ContainerImages {
		if _, err := os.Stat(in.Env.DockerFiles[i]); os.IsNotExist(err) {
			// path/to/whatever does not exist
			log.Info("skipping build because no dockerfile", log.M{"env": in.Env.DockerFiles[i]})
			continue
		}

		params := []string{"build", "-t", in.Env.ContainerImages[i], "-f", in.Env.DockerFiles[i], in.Env.DockerContexts[i]}

		if _, err := cmd(in.Flags.Verbose, "docker", params...); err != nil {
			return errors.Wrap(err, "unable to build docker image")
		}
	}

	return nil
}

func Push(in CommandInput) error {
	for i := range in.Env.ContainerImages {
		dolog := in.Flags.Verbose

		if in.Flags.Env == DevEnv {
			log.Info("skipping push because of flag", log.M{"env": in.Flags.Env})
			return nil
		}

		if _, err := os.Stat(in.Env.DockerFiles[i]); os.IsNotExist(err) {
			log.Info("skipping push because no dockerfile", log.M{"env": in.Env.DockerFiles[i]})
			continue
		}

		if _, err := cmd(dolog, in.Cloud, "docker", "--", "push", in.Env.ContainerImages[i]); err != nil {
			return errors.Wrap(err, "unable to push docker image")
		}
	}

	return nil
}
