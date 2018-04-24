package kubetool

import (
	"strings"

	"github.com/pkg/errors"
)

type Args struct {
	Components []string
}

type Flags struct {
	Verbose bool
	Env     string
}

type Env struct {
	Cloud             string            `envconfig:"CLOUD" default:"gcloud"`
	ContainerImages   []string          `envconfig:"CONTAINER_IMAGES" required:"true"`
	HelmChartPath     string            `envconfig:"HELM_CHART_PATH" required:"true"`
	HelmBaseValueFile string            `envconfig:"HELM_BASE_VALUE_FILE" required:"true"`
	HelmEnvValueFile  string            `envconfig:"HELM_ENV_VALUE_FILE" required:"true"`
	HelmImages        []string          `envconfig:"HELM_IMAGES" required:"true"`
	DockerFiles       []string          `envconfig:"DOCKER_FILES" required:"true"`
	DockerContexts    []string          `envconfig:"DOCKER_CONTEXTS" required:"true"`
	KubeContextMap    map[string]string `envconfig:"KUBE_CONTEXT_MAP" required:"true"`
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

	ciN := len(in.Env.ContainerImages)
	hiN := len(in.Env.HelmImages)
	dfN := len(in.Env.DockerFiles)
	dcN := len(in.Env.DockerContexts)
	if !((ciN == hiN) && (hiN == dfN) && (dfN == dcN)) {
		return errors.New("len(helm container images), len(helm images), len(docker files), len(docker contexts) must be equal")
	}
	return nil
}
