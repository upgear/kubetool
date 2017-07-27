package kubetool

import "github.com/pkg/errors"

func Push(in Input) error {
	dolog := in.Flags.Verbose

	if _, err := cmd(dolog, in.Env.Cloud, "docker", "--", "push", in.Env.ContainerImage); err != nil {
		return errors.Wrap(err, "unable to push docker image")
	}

	return nil
}
