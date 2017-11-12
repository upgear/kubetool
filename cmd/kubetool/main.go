package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/upgear/go-kit/log"
	"github.com/upgear/kubetool"
)

func main() {
	input := kubetool.Input{
		Env: kubetool.Env{
			Cloud:              envElse("KT_CLOUD", "gcloud"),
			ContainerImage:     envElse("KT_CONTAINER_IMAGE", "{{.Args.Name}}"),
			KubernetesFile:     envElse("KT_KUBERNETES_FILE", "."),
			KubernetesTestFile: envElse("KT_KUBERNETES_TEST_FILE", "."),
			DockerFile:         envElse("KT_DOCKER_FILE", "."),
			DockerContext:      envElse("KT_DOCKER_CONTEXT", "."),
		},
		Repo: kubetool.Repo{
			// Default to 'latest'
			CommitHash: "latest",
		},
	}

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage(""))
	}
	flag.BoolVar(&input.Flags.Verbose, "v", false, "Log a bunch of stuff")
	flag.BoolVar(&input.Flags.Latest, "latest", false, "Use 'latest' for .Repo.CommitHash in template")
	flag.BoolVar(&input.Flags.Save, "save", false, "Save deployed kubernetes config")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fatal(errors.New(usage("expected at least 2 arguments")))
	}
	input.Args = kubetool.Args{
		Commands: strings.Split(args[0], ","),
		Names:    args[1:],
	}

	if !input.Flags.Latest {
		fatal(kubetool.CheckRepo(&input))
	}

	// Map commands
	cmds := make([]kubetool.Command, len(input.Commands))
	for i, c := range input.Commands {
		var ok bool
		cmds[i], ok = kubetool.CommandMap[c]
		if !ok {
			fatal(fmt.Errorf("invalid command: %s", c))
		}
	}

	for _, input.Args.Name = range input.Args.Names {
		// Run commands
		for i, cmd := range cmds {
			if input.Flags.Verbose {
				log.Info("starting kubetool sub-command", log.M{
					"name":   input.Args.Name,
					"subcmd": input.Commands[i],
				})
			}
			fatal(input.Process())
			fatal(cmd(input))
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
		`Usage: kubetool [Options...] <command> <name>...

Commands:
    build   Runs docker build

    push    Runs docker push

    deploy  Runs docker build & kubectl apply

        Options:

            --save  Save the updated kubernetes config (with new image versions)

    test    Runs docker build & kubectl apply with test config

        Options:

            --save  (see deploy command)

Options:
    -h --help  Print usage
    -v         Verbose output

    --latest   Use 'latest' for .Repo.CommitHash in template

Environment Variables:
    KT_DOCKER_FILE      Dockerfile
    KT_KUBERNETES_FILE  Kubernetes config
    KT_CONTAINER_IMAGE  Template for container image (docker tag)
    KT_DOCKER_CONTEXT   Docker build context (directory)
    KT_CLOUD            Cloud provider (only supports 'gcloud')

Note:

    Commands can be comma seperated: blueprint build,push,deploy example

    Multiple names can be supplied:  blueprint build one two

    Multiple commands & names:       blueprint build,push,deploy one two
`
}
