package kubetool

import (
	"fmt"
	"strings"
)

type Command func(CommandInput) error

var CommandMap = map[string]Command{
	"build":    Build,
	"push":     Push,
	"undeploy": Undeploy,
	"deploy":   Deploy,
}

func ParseComponent(s string) (Component, error) {
	splt := strings.Split(s, "/")
	if len(splt) != 2 {
		return Component{}, fmt.Errorf("unable parse component %q into <domain>/<name>", s)
	}
	return Component{
		Domain: splt[0],
		Name:   splt[1],
	}, nil
}

// Component is made up of a domain and a component name.
type Component struct {
	Domain string
	Name   string
}

// ComandInput is the data that is passed to a command. It includes processed
// versions of environment templates.
type CommandInput struct {
	Component Component
	Env
	Flags
	Repo
}

// GetCommandInput translates RawInput into CommandInput for a given component
// by the component index (index must exist or panics).
func GetCommandInput(in RawInput, cmpIdx int) (cd CommandInput, err error) {
	// Parse component.
	cd.Component, err = ParseComponent(in.Args.Components[cmpIdx])
	if err != nil {
		return
	}

	// Parse env templates.
	tmplData := struct {
		Component
		Repo
	}{cd.Component, in.Repo}
	cd.Env.Cloud, err = templateString(in.Env.Cloud, tmplData)
	if err != nil {
		return
	}
	cd.Env.ContainerImage, err = templateString(in.Env.ContainerImage, tmplData)
	if err != nil {
		return
	}
	cd.Env.KubernetesFile, err = templateString(in.Env.KubernetesFile, tmplData)
	if err != nil {
		return
	}
	cd.Env.DockerFile, err = templateString(in.Env.DockerFile, tmplData)
	if err != nil {
		return
	}
	cd.Env.DockerContext, err = templateString(in.Env.DockerContext, tmplData)
	if err != nil {
		return
	}

	// Pass-thru data.
	cd.Flags = in.Flags
	cd.Repo = in.Repo

	return
}
