package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/upgear/go-kit/log"
	"github.com/upgear/kubetool"
)

func main() {
	// Set defaults.
	input := kubetool.RawInput{
		Repo: kubetool.Repo{
			Commit: "latest",
		},
	}

	envconfig.MustProcess("KT", &input.Env)

	// Parse flags.
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage(""))
	}
	flag.BoolVar(&input.Flags.Verbose, "v", false, "Log a bunch of stuff")
	flag.StringVar(&input.Flags.Env, "env", "dev", "Environment")
	var overrideRepoCheck bool
	flag.BoolVar(&overrideRepoCheck, "norepocheck", false, "Override repo integrity check")
	flag.Parse()

	// Parse arguments.
	args := flag.Args()
	if len(args) < 1 {
		fatal(errors.New(usage("expected at least 1 argument")))
	}
	input.Args = kubetool.Args{
		Components: args,
	}

	// Validate inputs.
	fatal(input.Process())

	if input.Flags.Verbose {
		log.GlobalLevel = log.LevelDebug
	} else {
		log.GlobalLevel = log.LevelError
	}

	// Inspect the repo.
	if !overrideRepoCheck {
		if err := kubetool.CheckRepo(&input); err != nil {
			if input.Flags.Env != kubetool.DevEnv {
				fatal(errors.Wrapf(err, "repo not acceptable for %q environment", input.Flags.Env))
			}
		}
	}

	// Set environment-specific values.
	switch input.Flags.Env {
	case kubetool.DevEnv:
		log.Info("setting docker env")
		fatal(kubetool.SetDevDockerEnv())
		// Set commit to a random string so that kubernetes will refresh every time.
		input.Repo.Commit = uuid.New().String()
	}

	// For each component.
	var state kubetool.State
	for cidx := range input.Args.Components {
		// Run all commands.
		for i, cmd := range kubetool.Commands {
			cmdInput, err := kubetool.GetCommandInput(input, cidx)
			fatal(errors.Wrap(err, "unable to parse input for component"))

			log.Info("starting kubetool sub-command", log.M{
				"chart":   cmdInput.Component.Chart,
				"release": cmdInput.Component.Release,
				"command": kubetool.CommandNames[i],
			})

			fatal(cmd(cmdInput, &state))
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
		`Usage: kubetool [Options...] <chart>/<release>...

Options:
    -e --env   Environment (defaults to "dev"). If it is anything other than
               "dev" then the repo must be in a clean state.
    -h --help  Print usage
    -v         Verbose output

Environment Variables:
    KT_COMPONENTS            Top level components (go, web, etc.)
    KT_DOCKER_FILE_DIRS      Directories that contain Dockerfiles
    KT_DOCKER_CONTEXTS       Docker build context (directories)

    KT_HELM_CHART_PATH       Helm chart (directory)
    KT_HELM_BASE_VALUE_FILE  The first layer of helm values
    KT_HELM_ENV_VALUE_FILE   The second layer of helm values (env specific)
    KT_DOCKER_REGISTRY_BASE  Template for docker registry
    KT_CLOUD                 Cloud provider (only supports 'gcloud')

    All environment variables can be templated using the following variables:

        {{.Env}}      Environment (i.e. "dev", "stg", "prd")
        {{.Chart}}    Helm chart
        {{.Release}}  Helm release name
        {{.Commit}}   Repo commit (git commit hash)

    Example: export KT_DOCKER_FILE="$GOPATH/src/{{.Chart}}/Dockerfile"

Note:

    Multiple names can be supplied:  kubetool foo/bar abc/xyz
`
}
