package kubetool

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Command func(CommandInput, *State) error

const DevEnv = "dev"

// TODO: Revisit "helm test command"
var Commands = []Command{Build, Push, Apply}
var CommandNames = []string{"build", "push", "apply"}

func ParseComponent(s string) (Component, error) {
	splt := strings.Split(s, "/")
	if len(splt) != 2 {
		return Component{}, fmt.Errorf("unable parse component %q into <chart>/<release>", s)
	}
	return Component{
		Chart:   splt[0],
		Release: splt[1],
	}, nil
}

// Component is made up of a helm chart and a release name.
type Component struct {
	Chart   string
	Release string
}

func (c Component) ChartRelease() string {
	return fmt.Sprintf("%s-%s", c.Chart, c.Release)
}

// ComandInput is the data that is passed to a command. It includes processed
// versions of environment templates.
type CommandInput struct {
	Component
	Env
	Flags
	Repo
}

// EnvTemplateData is the data used to excute environment variable templating.
type EnvTemplateData struct {
	Component
	Repo
	// Env is the software environment (ie. "dev", "stg", etc.)
	Env string
}

// KubeTemplateData is the data used to execute a kubernetes template file.
type KubeTemplateData struct {
}

// GetCommandInput translates RawInput into CommandInput for a given component
// by the component index (index must exist or panics).
func GetCommandInput(in RawInput, cmpIdx int) (cd CommandInput, err error) {
	// Parse component.
	cd.Component, err = ParseComponent(in.Args.Components[cmpIdx])
	if err != nil {
		return
	}

	// Consider getting rid of this field copying b/c its brittle to change.
	cd.Env.Components = in.Env.Components
	cd.Env.DockerFileDirs = make([]string, len(in.Env.DockerFileDirs))
	cd.Env.DockerContexts = make([]string, len(in.Env.DockerContexts))
	cd.Env.KubeContextMap = in.Env.KubeContextMap
	cd.Env.DockerRegistryBase = in.Env.DockerRegistryBase

	// Parse env templates.
	tmplData := EnvTemplateData{cd.Component, in.Repo, in.Flags.Env}
	cd.Env.Cloud, err = templateString(in.Env.Cloud, tmplData)
	if err != nil {
		return
	}
	cd.Env.HelmChartPath, err = templateString(in.Env.HelmChartPath, tmplData)
	if err != nil {
		return
	}
	cd.Env.HelmBaseValueFile, err = templateString(in.Env.HelmBaseValueFile, tmplData)
	if err != nil {
		return
	}
	cd.Env.HelmEnvValueFile, err = templateString(in.Env.HelmEnvValueFile, tmplData)
	if err != nil {
		return
	}
	for i := range in.Env.DockerFileDirs {
		cd.Env.DockerFileDirs[i], err = templateString(in.Env.DockerFileDirs[i], tmplData)
		if err != nil {
			return
		}
	}
	for i := range in.Env.DockerContexts {
		cd.Env.DockerContexts[i], err = templateString(in.Env.DockerContexts[i], tmplData)
		if err != nil {
			return
		}
	}

	// Pass-thru data.
	cd.Flags = in.Flags
	cd.Repo = in.Repo

	return
}

type Args struct {
	Components []string
}

type Flags struct {
	Verbose bool
	Env     string
}

type Env struct {
	Cloud             string            `envconfig:"CLOUD" default:"gcloud"`
	HelmChartPath     string            `envconfig:"HELM_CHART_PATH" required:"true"`
	HelmBaseValueFile string            `envconfig:"HELM_BASE_VALUE_FILE" required:"true"`
	HelmEnvValueFile  string            `envconfig:"HELM_ENV_VALUE_FILE" required:"true"`
	KubeContextMap    map[string]string `envconfig:"KUBE_CONTEXT_MAP" required:"true"`

	DockerRegistryBase string `envconfig:"DOCKER_REGISTRY_BASE" required:"true"`

	Components     []string `envconfig:"COMPONENTS" required:"true"`
	DockerContexts []string `envconfig:"DOCKER_CONTEXTS" required:"true"`
	DockerFileDirs []string `envconfig:"DOCKER_FILE_DIRS" required:"true"`
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

	hiN := len(in.Env.Components)
	dfN := len(in.Env.DockerFileDirs)
	dcN := len(in.Env.DockerContexts)
	if !((hiN == dfN) && (dfN == dcN)) {
		return errors.New("len(helm container images), len(helm images), len(docker files), len(docker contexts) must be equal")
	}
	return nil
}

type State struct {
	DockerTags []DockerTag
}

type DockerTag struct {
	Key string
	Tag string
}
