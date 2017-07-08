# Kubetool

A tool for working with kubernetes.

## Example Usage

```sh
# In a repo with "example.Dockerfile"

export GCLOUD_PROJECT_ID=abc
export KT_DOCKER_TAG="gcr.io/$GCLOUD_PROJECT_ID/{{.Args.Name}}:{{.Repo.CommitHash}}"

kubetool -v build example
```
