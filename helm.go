package kubetool

import "os"

func Install(in CommandInput) (err error) {
	args := []string{
		"install", in.HelmChartPath,
		"--name", in.Release,
		"--values", in.HelmBaseValueFile,
	}

	if err := appendEnvConfig(in, &args); err != nil {
		return err
	}

	_, err = cmd(in.Flags.Verbose, "helm", args...)
	return
}

func Upgrade(in CommandInput) (err error) {
	args := []string{
		"upgrade", in.Release, in.HelmChartPath,
		"--values", in.HelmBaseValueFile,
	}

	if err := appendEnvConfig(in, &args); err != nil {
		return err
	}

	_, err = cmd(in.Flags.Verbose, "helm", args...)

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
