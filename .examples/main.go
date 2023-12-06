package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/clickcode-dev/fiberdi"
	"github.com/clickcode-dev/fiberdi/.examples/cat"
	"github.com/clickcode-dev/fiberdi/.examples/dog"
	"github.com/gofiber/fiber/v2"
)

func main() {
	appModule := &fiberdi.Module{
		Modules: []fiberdi.IModule{
			dog.Module,
			cat.Module,
		},
	}

	app := fiberdi.New(appModule, fiber.Config{})

	port := "3000"

	log.Infof("Server is running on port %s", port)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
