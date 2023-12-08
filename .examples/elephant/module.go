package elephant

import (
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/.examples/dog"
)

var Module = &fiberdi.Module{
	Controllers: []fiberdi.IController{
		&ElephantController{},
	},
	Injectables: []interface{}{
		&ElephantService{},
		&dog.DogService{},
	},
	Imports: []fiberdi.IModule{},
}
