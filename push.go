package kubetool

import "github.com/pkg/errors"

func Push(in Input) error {
	dolog := in.Flags.Verbose

	tag, err := templateString(in.Env.TagTemplate, in)
	if err != nil {
		return errors.Wrap(err, "unable to template tag")
	}

	if _, err := cmd(dolog, in.Env.Cloud, "docker", "--", "push", tag); err != nil {
		return errors.Wrap(err, "unable to push docker image")
	}

	return nil
}
