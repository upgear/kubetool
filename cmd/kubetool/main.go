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
			KubernetesPath: envElse("KT_KUBERNETES_PATH", "."),
			DockerfilePath: envElse("KT_DOCKERFILE_PATH", "."),
			DockerContext:  envElse("KT_DOCKER_CONTEXT", "."),
			TagTemplate:    envElse("KT_DOCKER_TAG", "{{.Args.Name}}"),
		},
	}

	flag.BoolVar(&input.Flags.Verbose, "v", false, "Log a bunch of stuff")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fatal(errors.New(usage("expected exactly 2 arguments")))
	}
	input.Args = kubetool.Args{
		Command: args[0],
		Name:    args[1],
	}

	fatal(kubetool.CheckRepo(&input))

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
	return msg + "\n\n" +
		`Usage: kubetool [Options...] <command> <name>

Commands:
    build   Runs docker build
    deploy  Runs docker build & kubectl apply

Options:
    -v      Verbose output

Environment Variables:
    KT_DOCKERFILE_PATH   Directory to look for Dockerfiles
    KT_KUBERNETES_PATH   Direcotry to look for kubernetes configs
    KT_DOCKER_TAG        Template to create docker tag
    KT_DOCKER_CONTEXT    Docker build context (directory)
`
}
