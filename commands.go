package kubetool

import (
	"fmt"
	"strings"
)

type Command func(CommandInput) error

const DevEnv = "dev"

var CommandMap = map[string]Command{
	"build":  Build,
	"push":   Push,
	"apply":  Apply,
	"test":   Test,
	"kill":   Kill,
	"delete": Delete,
}

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

	cd.Env.HelmImages = make([]string, len(in.Env.HelmImages))
	cd.Env.ContainerImages = make([]string, len(in.Env.ContainerImages))
	cd.Env.DockerFiles = make([]string, len(in.Env.DockerFiles))
	cd.Env.DockerContexts = make([]string, len(in.Env.DockerContexts))

	// Parse env templates.
	tmplData := EnvTemplateData{cd.Component, in.Repo, in.Flags.Env}
	cd.Env.Cloud, err = templateString(in.Env.Cloud, tmplData)
	if err != nil {
		return
	}
	for i := range in.Env.HelmImages {
		cd.Env.HelmImages[i], err = templateString(in.Env.HelmImages[i], tmplData)
		if err != nil {
			return
		}
	}
	for i := range in.Env.ContainerImages {
		cd.Env.ContainerImages[i], err = templateString(in.Env.ContainerImages[i], tmplData)
		if in.Flags.Env == DevEnv {
			cd.Env.ContainerImages[i] = strings.Replace(cd.Env.ContainerImages[i], "/", "-", -1)
		}
		if err != nil {
			return
		}
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
	for i := range in.Env.DockerFiles {
		cd.Env.DockerFiles[i], err = templateString(in.Env.DockerFiles[i], tmplData)
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
