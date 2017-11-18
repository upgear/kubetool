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
			Cloud:          envElse("KT_CLOUD", "gcloud"),
			ContainerImage: envElse("KT_CONTAINER_IMAGE", "{{.Name}}"),
			KubernetesFile: envElse("KT_KUBERNETES_FILE", "."),
			DockerFile:     envElse("KT_DOCKER_FILE", "."),
			DockerContext:  envElse("KT_DOCKER_CONTEXT", "."),
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
	flag.BoolVar(&input.Flags.Local, "local", false, "Use 'latest' for .Repo.CommitHash in template and skip push")
	flag.BoolVar(&input.Flags.Save, "save", false, "Save templated kubernetes config from deploy")
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

	// Inspect the repo.
	if !input.Flags.Local {
		fatal(kubetool.CheckRepo(&input))
	}

	// Validate inputs.
	fatal(input.Validate())

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
					"domain":    cmdInput.Component.Domain,
					"component": cmdInput.Component.Name,
					"subcmd":    input.Args.Commands[i],
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
		`Usage: kubetool [Options...] <command> <domain>/<name>...

Commands:
    build     Runs docker build

    push      Runs docker push

	undeploy  Delete previous kubernetes deployments

    deploy    Runs kubectl apply

        Options:

            --save  Save the updated kubernetes config (with new image versions)

Options:
    -h --help  Print usage
    -v         Verbose output

    --local    Use 'latest' for .CommitHash in env template and skip pushes

Environment Variables:
    KT_DOCKER_FILE      Dockerfile
    KT_KUBERNETES_FILE  Kubernetes config
    KT_CONTAINER_IMAGE  Template for container image (docker tag)
    KT_DOCKER_CONTEXT   Docker build context (directory)
    KT_CLOUD            Cloud provider (only supports 'gcloud')

    All environment variables can be templated using the following variables:

        {{.Domain}}  Component domain
        {{.Name}}    Component name
        {{.Commit}}  Repo commit (git commit hash)

    Example: export KT_DOCKER_FILE="$GOPATH/src/{{.Domain}}/{{.Name}}/Dockerfile"

Note:

    Commands can be comma seperated: kubetool build,push,deploy foo/bar

    Multiple names can be supplied:  kubetool build foo/bar abc/xyz

    Multiple commands & names:       kubetool build,push,deploy foo/bar abc/xyz
`
}
