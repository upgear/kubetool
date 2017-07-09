package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/upgear/kubetool"
)

func main() {
	input := kubetool.Input{
		Env: kubetool.Env{
			Cloud:          envElse("KT_CLOUD", "gcloud"),
			ContainerImage: envElse("KT_CONTAINER_IMAGE", "{{.Args.Name}}"),
			KubernetesPath: envElse("KT_KUBERNETES_PATH", "."),
			DockerfilePath: envElse("KT_DOCKERFILE_PATH", "."),
			DockerContext:  envElse("KT_DOCKER_CONTEXT", "."),
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
	flag.BoolVar(&input.Flags.Latest, "latest", false, "Use 'latest' as container image version rather than git commit")
	flag.BoolVar(&input.Flags.Save, "save", false, "Save deployed kubernetes config")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fatal(errors.New(usage("expected exactly 2 arguments")))
	}
	input.Args = kubetool.Args{
		Command: args[0],
		Name:    args[1],
	}

	fatal(input.Validate())
	if !input.Flags.Latest {
		fatal(kubetool.CheckRepo(&input))
	}

	cmd, ok := kubetool.Commands[input.Command]
	if !ok {
		fatal(fmt.Errorf("invalid command: %s", input.Command))
	}

	fatal(cmd(input))
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
		`Usage: kubetool [Options...] <command> <name>

Commands:
    build   Runs docker build

    push    Runs docker push

    deploy  Runs docker build & kubectl apply

        Options:

            --save  Save the updated kubernetes config (with new image versions)

Options:
    -h --help  Print usage
    -v         Verbose output

    --latest   Use 'latest' as container image version rather than git commit

Environment Variables:
    KT_DOCKERFILE_PATH  Directory to look for Dockerfiles
    KT_KUBERNETES_PATH  Direcotry to look for kubernetes configs
    KT_CONTAINER_IMAGE  Template for container image (docker tag)
    KT_DOCKER_CONTEXT   Docker build context (directory)
    KT_CLOUD            Cloud provider (only supports 'gcloud')
`
}
