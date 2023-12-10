package dog

import (
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/.examples/cat"
)

var Module = &fiberdi.Module{
	Controllers: []fiberdi.IController{
		new(DogController),
	},
	Injectables: []interface{}{
		new(DogService),
		new(cat.CatService),
	},
	Imports: []fiberdi.IModule{
		cat.Module,
	},
}
