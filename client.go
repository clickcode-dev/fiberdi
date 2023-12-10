package fiberdi

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func New(module IModule, configs ...fiber.Config) *fiber.App {
	config := fiber.Config{}

	if len(configs) > 0 {
		config = configs[0]
	}

	config.ReduceMemoryUsage = true
	config.StrictRouting = true
	config.CaseSensitive = true
	config.DisableStartupMessage = true
	config.JSONEncoder = sonic.Marshal
	config.JSONDecoder = sonic.Unmarshal

	app := fiber.New(config)

	module.setInjectableCurrentInAppModuleFn(module.(*Module).setInjectableCurrent)
	module.setInjectablesNames(module.mappedInjectables(module.(*Module).Modules, []string{}))
	module.setExportsNames(module.mappedExports(module.(*Module).Modules, []string{}))

	app = module.start(app)

	return app
}
