package kubetool

func Install(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "install", in.HelmChartPath,
		"--name", in.Release,
		"--values", in.HelmBaseValueFile,
		"--values", in.HelmEnvValueFile,
	)
	return
}

func Upgrade(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "upgrade", in.Release, in.HelmChartPath,
		"--values", in.HelmBaseValueFile,
		"--values", in.HelmEnvValueFile,
	)
	return
}

func Delete(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "delete", in.Component.Release)
	return
}
