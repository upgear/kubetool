package kubetool

type Input struct {
	Args
	Flags
	Env
	Repo
}

type Args struct {
	Command string
	Name    string
}

type Flags struct {
	Verbose bool
}

type Env struct {
	TagTemplate    string
	KubernetesPath string
	DockerfilePath string
	DockerContext  string
}

type Repo struct {
	CommitHash string
}
