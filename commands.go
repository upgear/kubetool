package kubetool

type Command func(Input) error

var Commands = map[string]Command{
	"build":  Build,
	"deploy": Deploy,
}
