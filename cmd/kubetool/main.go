package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/upgear/go-kit/log"
	"github.com/upgear/kubetool"
)

func main() {
	// Set defaults.
	input := kubetool.RawInput{
		Env: kubetool.Env{
			Cloud:             envElse("KT_CLOUD", "gcloud"),
			ContainerImage:    envElse("KT_CONTAINER_IMAGE", "{{.Name}}"),
			HelmChartPath:     envElse("KT_HELM_CHART_PATH", "."),
			HelmBaseValueFile: envElse("KT_HELM_BASE_VALUE_FILE", "."),
			HelmEnvValueFile:  envElse("KT_HELM_ENV_VALUE_FILE", "."),
			DockerFile:        envElse("KT_DOCKER_FILE", "."),
			DockerContext:     envElse("KT_DOCKER_CONTEXT", "."),
		},
		Repo: kubetool.Repo{
			Commit: "latest",
		},
	}

	// Parse flags.
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage(""))
	}
	flag.BoolVar(&input.Flags.Verbose, "v", false, "Log a bunch of stuff")
	flag.StringVar(&input.Flags.Env, "env", "dev", "Environment")
	flag.Parse()

	// Parse arguments.
	args := flag.Args()
	if len(args) < 2 {
		fatal(errors.New(usage("expected at least 2 arguments")))
	}
	input.Args = kubetool.Args{
		Commands:   strings.Split(args[0], ","),
		Components: args[1:],
	}

	// Validate inputs.
	fatal(input.Process())

	// Inspect the repo.
	if err := kubetool.CheckRepo(&input); err != nil {
		if input.Flags.Env != kubetool.DevEnv {
			fatal(errors.Wrapf(err, "repo not acceptable for %q environment", input.Flags.Env))
		}
	}

	// Map commands.
	cmds := make([]kubetool.Command, len(input.Args.Commands))
	for i, c := range input.Args.Commands {
		var ok bool
		cmds[i], ok = kubetool.CommandMap[c]
		if !ok {
			fatal(fmt.Errorf("invalid command: %s", c))
		}
	}

	// For each component.
	for cidx := range input.Args.Components {
		// Run all commands.
		for i, cmd := range cmds {
			cmdInput, err := kubetool.GetCommandInput(input, cidx)
			fatal(errors.Wrap(err, "unable to parse input for component"))

			if input.Flags.Verbose {
				log.Info("starting kubetool sub-command", log.M{
					"chart":   cmdInput.Component.Chart,
					"release": cmdInput.Component.Release,
					"command": input.Args.Commands[i],
				})
			}

			fatal(cmd(cmdInput))
		}
	}
}

func envElse(env, def string) string {
	s, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	return s
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func usage(msg string) string {
	if msg != "" {
		msg = msg + "\n\n"
	}
	return msg +
		`Usage: kubetool [Options...] <comma-seperated-commands> <chart>/<release>...

Commands:
    build     Runs docker build

    push      Runs docker push

    install   Runs helm install

    upgrade   Runs helm upgrade

    test      Runs helm test

    kill      Runs kubectl delete pod (useful for development)

    delete    Runs helm delete

Options:
    -e --env   Environment (defaults to "dev"). If it is anything other than
               "dev" then the repo must be in a clean state.
    -h --help  Print usage
    -v         Verbose output

Environment Variables:
    KT_DOCKER_FILE           Dockerfile
    KT_HELM_CHART_PATH       Helm chart (directory)
    KT_HELM_BASE_VALUE_FILE  The first layer of helm values
    KT_HELM_ENV_VALUE_FILE   The second layer of helm values (env specific)
    KT_CONTAINER_IMAGE       Template for container image (docker tag)
    KT_DOCKER_CONTEXT        Docker build context (directory)
    KT_CLOUD                 Cloud provider (only supports 'gcloud')

    All environment variables can be templated using the following variables:

        {{.Env}}      Environment (i.e. "dev", "stg", "prd")
        {{.Chart}}    Helm chart
        {{.Release}}  Helm release name
        {{.Commit}}   Repo commit (git commit hash)

    Example: export KT_DOCKER_FILE="$GOPATH/src/{{.Chart}}/Dockerfile"

Note:

    Commands can be comma seperated: kubetool build,push,install foo/bar

    Multiple names can be supplied:  kubetool build foo/bar abc/xyz

    Multiple commands & names:       kubetool build,push,install foo/bar abc/xyz
`
}
