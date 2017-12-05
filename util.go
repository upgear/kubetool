package kubetool

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func cmd(dolog bool, name string, params ...string) (string, error) {
	cmd := exec.Command(name, params...)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr

	if dolog {
		cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
		logCmd(name, params...)
	} else {
		cmd.Stdout = &stdout
	}

	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "unable to execute command: %s %s: %s", name, strings.Join(params, " "), stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func logCmd(name string, params ...string) {
	log.Info("running command", log.M{
		"cmd": fmt.Sprintf("%s %s", name, strings.Join(params, " ")),
	})
}

func templateString(s string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
