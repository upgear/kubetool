package kubetool

type Input struct {
	Args
	Flags
	Env
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
