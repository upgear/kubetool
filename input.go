package kubetool

import (
	"strings"

	"github.com/pkg/errors"
)

type Args struct {
	Commands   []string
	Components []string
}

type Flags struct {
	Verbose bool
	Env     string
}

type Env struct {
	Cloud             string
	ContainerImage    string
	HelmChartPath     string
	HelmBaseValueFile string
	HelmEnvValueFile  string
	DockerFile        string
	DockerContext     string
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

func (in *RawInput) Process() error {
	in.Flags.Env = strings.ToLower(in.Flags.Env)

	if in.Env.Cloud != "gcloud" {
		return errors.New("only 'gcloud' is a supported cloud type")
	}
	return nil
}
