package cat

import (
	"github.com/clickcode-dev/fiberdi"
)

var Module = &fiberdi.Module{
	Controllers: []fiberdi.IController{
		&CatController{},
	},
	Injectables: []interface{}{
		&CatService{},
	},
}
