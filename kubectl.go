package kubetool

func Kill(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose,
		"kubectl",
		"--context="+kubeContext(in),
		"delete", "pod",
		"-l", "release="+in.Component.ChartRelease(),
	)

	return
}

func kubeContext(in CommandInput) string {
	ctx, ok := in.Env.KubeContextMap[in.Flags.Env]
	if !ok {
		return "minikube"
	}
	return ctx
}
