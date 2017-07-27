package kubetool

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func Deploy(in Input) error {
	file := in.ComputedEnv.KubernetesFile

	confBtys, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "unable to read kubernetes file: %s", file)
	}

	splitTag := strings.Split(in.ComputedEnv.ContainerImage, ":")
	if len(splitTag) != 2 {
		return errors.New("expected tag to have version")
	}

	newConf, err := updateImage(confBtys, splitTag[0], splitTag[1])
	if err != nil {
		return errors.Wrap(err, "unable to update image version")
	}

	name := "kubectl"
	params := []string{"apply", "-f", "-"}
	cmd := exec.Command(name, params...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "unable to open stdin to kubectl")
	}

	go func() {
		stdin.Write(newConf)
		stdin.Close()
	}()

	if in.Flags.Verbose {
		cmd.Stdout = os.Stdout
		log.Info("modified tag", log.M{"tag": in.ComputedEnv.ContainerImage})
		fmt.Println(string(newConf))
		logCmd(name, params...)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "unable to execute command: %s %s: %s", name, strings.Join(params, " "), stderr.String())
	}

	if in.Flags.Save {
		if err := ioutil.WriteFile(file, newConf, 0644); err != nil {
			return errors.Wrapf(err, "unable to write updated kubernetes config file: %s", file)
		}
	}

	return nil
}

func updateImage(conf []byte, pretag, version string) ([]byte, error) {
	rgx, err := regexp.Compile(fmt.Sprintf(`image:\s+%s:(\S+)`, pretag))
	if err != nil {
		return nil, err
	}

	return rgx.ReplaceAll(conf, []byte(fmt.Sprintf("image: %s:%s", pretag, version))), nil
}
