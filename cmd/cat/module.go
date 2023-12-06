package cat

import (
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/cmd/dog"
)

var Module = fiberdi.Module{
	Controllers: []fiberdi.IController{
		&CatController{},
	},
	Injectables: []interface{}{
		&CatService{},
		&dog.DogService{},
	},
}
