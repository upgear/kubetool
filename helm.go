package kubetool

import (
	"os"
	"strings"
)

func Install(in CommandInput) (err error) {
	args := []string{
		"--kube-context", kubeContext(in),
		"install", in.HelmChartPath,
		"--name", in.ChartRelease(),
		"--values", in.HelmBaseValueFile,
		"--set", helmVals(in),
	}
	if in.Verbose {
		args = append(args, "--debug")
	}

	if err := appendEnvConfig(in, &args); err != nil {
		return err
	}

	_, err = cmd(in.Flags.Verbose, "helm", args...)
	return
}

func helmVals(in CommandInput) string {
	var imgs []string
	for i, img := range in.HelmImages {
		imgs = append(imgs, img+"="+in.ContainerImages[i])
	}
	return strings.Join(imgs, ",")
}

func Upgrade(in CommandInput) (err error) {
	args := []string{
		"--kube-context", kubeContext(in),
		"upgrade", in.ChartRelease(), in.HelmChartPath,
		"--values", in.HelmBaseValueFile,
		"--set", helmVals(in),
	}
	if in.Verbose {
		args = append(args, "--debug")
	}

	if err := appendEnvConfig(in, &args); err != nil {
		return err
	}

	_, err = cmd(in.Flags.Verbose, "helm", args...)

	return
}

func Test(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "--kube-context", kubeContext(in), "test", "--debug", "--cleanup", in.Component.ChartRelease())

	return
}

func Delete(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "--kube-context", kubeContext(in), "delete", in.Component.ChartRelease())
	return
}

// appendEnvConfig is the file exists.
func appendEnvConfig(in CommandInput, args *[]string) error {
	_, err := os.Stat(in.HelmEnvValueFile)

	if os.IsNotExist(err) {
		// The env config file does not exist.
		return nil
	}

	if err == nil {
		// The env config file exists.
		*args = append(*args, "--values", in.HelmEnvValueFile)
	}

	// Something went wrong when trying to stat.
	return err
}
