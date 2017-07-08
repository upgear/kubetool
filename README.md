# Kubetool

A tool for working with kubernetes.

## Example Usage

```sh
# In a repo containing example.Dockerfile and example.yaml ...

export GCLOUD_PROJECT_ID=abc
export KT_DOCKER_TAG="gcr.io/$GCLOUD_PROJECT_ID/{{.Args.Name}}:{{.Repo.CommitHash}}"

kubetool build example
kubetool deploy example
```
