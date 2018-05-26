package kubetool

import (
	"os"
	"strings"
)

func helmVals(in CommandInput, s *State) string {
	var imgs []string
	for _, tag := range s.DockerTags {
		imgs = append(imgs, tag.Key+"_image="+tag.Tag)
	}
	return strings.Join(imgs, ",")
}

func Apply(in CommandInput, s *State) (err error) {
	args := []string{
		"--kube-context", kubeContext(in),
		"upgrade", in.ChartRelease(), in.HelmChartPath,
		"--install",
		"--values", in.HelmBaseValueFile,
		"--set", helmVals(in, s),
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

/*
// TODO: Revisit
func Test(in CommandInput) (err error) {
	_, err = cmd(in.Flags.Verbose, "helm", "--kube-context", kubeContext(in), "test", "--debug", "--cleanup", in.Component.ChartRelease())

	return
}
*/

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
