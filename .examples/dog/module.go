package dog

import "github.com/clickcode-dev/fiberdi"

var Module = fiberdi.Module{
	Controllers: []fiberdi.IController{
		&DogController{},
	},
	Injectables: []interface{}{
		&DogService{},
	},
}
