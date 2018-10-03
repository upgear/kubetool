package kubetool

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/upgear/go-kit/log"
)

func Build(in CommandInput, s *State) error {
	for i, c := range in.Env.Components {
		files, err := ioutil.ReadDir(in.Env.DockerFileDirs[i])
		if err != nil {
			return errors.Wrap(err, "reading dir")
		}

		for _, f := range files {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".Dockerfile" {
				fp := filepath.Join(in.Env.DockerFileDirs[i], f.Name())
				dockerfileName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))

				tagLatest := buildTagName(in.Flags.Env == DevEnv, c, in.Env.DockerRegistryBase, dockerfileName, "latest")
				tag := buildTagName(in.Flags.Env == DevEnv, c, in.Env.DockerRegistryBase, dockerfileName, in.Repo.Commit)

				var params []string
				if in.OnlyBuild == "" || strings.Contains(f.Name(), in.OnlyBuild) {
					params = []string{"build", "-t", tagLatest, "-t", tag, "-f", fp, in.Env.DockerContexts[i]}
				} else {
					if in.Flags.Verbose {
						log.Info("excluding docker build, using latest tag", log.M{"dockerfile": f.Name()})
					}
					params = []string{"tag", tagLatest, tag}
				}

				if _, err := cmd(in.Flags.Verbose, "docker", params...); err != nil {
					return errors.Wrap(err, "unable to build/tag docker image")
				}

				s.DockerTags = append(s.DockerTags, DockerTag{
					Tag: tag,
					Key: c + "_" + dockerfileName,
				})
			}
		}
	}

	return nil
}

func buildTagName(dev bool, cmp, base, name, commit string) string {
	if dev && commit != "latest" {
		commit = commit[:8]
	}
	s := fmt.Sprintf("%s/%s-%s:%s", base, cmp, name, commit)
	if dev {
		s = strings.Replace(s, "/", "-", -1)
	}
	return s
}

func Push(in CommandInput, s *State) error {
	for _, tag := range s.DockerTags {
		dolog := in.Flags.Verbose

		if in.Flags.Env == DevEnv {
			log.Info("skipping push because of flag", log.M{"env": in.Flags.Env})
			return nil
		}

		if _, err := cmd(dolog, in.Cloud, "docker", "--", "push", tag.Tag); err != nil {
			return errors.Wrap(err, "unable to push docker image")
		}
	}

	return nil
}
