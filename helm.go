package kubetool

import (
	"os"
	"strings"
)

func Install(in CommandInput) (err error) {
	args := []string{
		"install", in.HelmChartPath,
		"--name", in.Release,
		"--values", in.HelmBaseValueFile,
		"--set", helmVals(in),
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
		"upgrade", in.Release, in.HelmChartPath,
		"--values", in.HelmBaseValueFile,
		"--set", helmVals(in),
	}

	if err := appendEnvConfig(in, &args); err != nil {
		return err
	}

	_, err = cmd(in.Flags.Verbose, "helm", args...)

	return
}

func Test(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "test", "--debug", "--cleanup", in.Component.Release)
	return
}

func Delete(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "delete", in.Component.Release)
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
