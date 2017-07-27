package kubetool

import "github.com/pkg/errors"

type Input struct {
	Args
	Flags
	Env
	Repo
}

func (in *Input) Process() error {
	if err := in.Env.Process(); err != nil {
		return errors.Wrap(err, "invalid environment")
	}

	var err error
	in.Env.Cloud, err = templateString(in.Env.Cloud, *in)
	if err != nil {
		return err
	}
	in.Env.ContainerImage, err = templateString(in.Env.ContainerImage, *in)
	if err != nil {
		return err
	}
	in.Env.KubernetesFile, err = templateString(in.Env.KubernetesFile, *in)
	if err != nil {
		return err
	}
	in.Env.DockerFile, err = templateString(in.Env.DockerFile, *in)
	if err != nil {
		return err
	}
	in.Env.DockerContext, err = templateString(in.Env.DockerContext, *in)
	if err != nil {
		return err
	}

	return nil
}

type Args struct {
	Commands []string
	Names    []string
	Name     string
}

type Flags struct {
	Verbose bool
	Latest  bool
	Save    bool
}

type Env struct {
	Cloud          string
	ContainerImage string
	KubernetesFile string
	DockerFile     string
	DockerContext  string
}

func (env Env) Process() error {
	if env.Cloud != "gcloud" {
		return errors.New("only 'gcloud' is a supported cloud type")
	}

	return nil
}

type Repo struct {
	CommitHash string
}
