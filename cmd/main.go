package main

import (
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/cmd/cat"
	"github.com/clickcode-dev/fiberdi/cmd/dog"
)

func main() {
	app := fiberdi.New(&fiberdi.Module{
		Modules: []fiberdi.Module{
			dog.Module,
			cat.Module,
		},
	})

	app.Listen(":3000")
}
