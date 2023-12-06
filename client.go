package fiberdi

import (
	"github.com/gofiber/fiber/v2"
)

func New(module *Module, config ...fiber.Config) *fiber.App {
	app := fiber.New(config...)

	module.foundInjectables(module, new([]string))
	module.foundExports(module, new([]string), new([]interface{}))

	app = module.addDependencies(app)

	return app
}
