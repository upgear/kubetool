package kubetool

type Command func(Input) error

var CommandMap = map[string]Command{
	"build":    Build,
	"push":     Push,
	"undeploy": Undeploy,
	"deploy":   Deploy,
}
