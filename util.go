package kubetool

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func cmd(dolog bool, name string, params ...string) (string, error) {
	cmd := exec.Command(name, params...)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if dolog {
		logCmd(name, params...)
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

func templateString(s string, in Input) (string, error) {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, in); err != nil {
		return "", err
	}

	return b.String(), nil
}
