package kubetool

import (
	"github.com/pkg/errors"
)

type Args struct {
	Commands   []string
	Components []string
}

type Flags struct {
	Verbose bool
	Local   bool
	Save    bool
}

type Env struct {
	Cloud          string
	ContainerImage string
	KubernetesFile string
	DockerFile     string
	DockerContext  string
}

type Repo struct {
	Commit string
}

type RawInput struct {
	Args  Args
	Flags Flags
	Env   Env
	Repo  Repo
}

func (in RawInput) Validate() error {
	if in.Env.Cloud != "gcloud" {
		return errors.New("only 'gcloud' is a supported cloud type")
	}
	return nil
}
