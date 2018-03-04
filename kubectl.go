package kubetool

func Kill(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "kubectl", "delete", "pod", "-l", "release="+in.Component.ChartRelease())
	return
}
