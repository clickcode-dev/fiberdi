package elephant

import (
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/.examples/dog"
)

var Module = &fiberdi.Module{
	Controllers: []fiberdi.IController{
		new(ElephantController),
	},
	Injectables: []interface{}{
		new(ElephantService),
		new(dog.DogService),
	},
	Imports: []fiberdi.IModule{
		dog.Module,
	},
}
