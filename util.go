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

func SetDevDockerEnv() error {
	env, err := minikubeDockerEnv()
	if err != nil {
		return err
	}
	setEnvs(parseEnvExports(env))
	return nil
}

func minikubeDockerEnv() (string, error) {
	c := exec.Command("minikube", "docker-env")
	var buf bytes.Buffer
	c.Stdout = &buf
	if err := c.Run(); err != nil {
		return "", errors.Wrap(err, "running command 'minikube docker-env'")
	}
	return buf.String(), nil
}

type kv struct {
	key string
	val string
}

func parseEnvExports(s string) []kv {
	lines := strings.Split(s, "\n")

	var kvs []kv

	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "export") {
			ln = strings.TrimSpace(strings.TrimPrefix(ln, "export"))
			splt := strings.Split(ln, "=")
			if len(splt) == 2 {
				kvs = append(kvs, kv{
					key: strings.TrimSpace(splt[0]),
					val: strings.TrimSpace(
						strings.TrimSuffix(strings.TrimPrefix(splt[1], `"`), `"`),
					),
				})
			}
		}
	}

	return kvs
}

func setEnvs(envs []kv) {
	for _, kv := range envs {
		os.Setenv(kv.key, kv.val)
	}
}
