package cat

import (
	"github.com/clickcode-dev/fiberdi"
)

var Module = &fiberdi.Module{
	Controllers: []fiberdi.IController{
		new(CatController),
	},
	Injectables: []interface{}{
		new(CatService),
	},
}
