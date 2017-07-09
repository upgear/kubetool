package kubetool

import "github.com/pkg/errors"

type Input struct {
	Args
	Flags
	Env
	Repo
}

func (in Input) Validate() error {
	if err := in.Env.Validate(); err != nil {
		return errors.Wrap(err, "invalid environment")
	}

	return nil
}

type Args struct {
	Command string
	Name    string
}

type Flags struct {
	Verbose bool
	Latest  bool
	Save    bool
}

type Env struct {
	Cloud          string
	ContainerImage string
	KubernetesPath string
	DockerfilePath string
	DockerContext  string
}

func (env Env) Validate() error {
	if env.Cloud != "gcloud" {
		return errors.New("only 'gcloud' is a supported cloud type")
	}

	return nil
}

type Repo struct {
	CommitHash string
}
