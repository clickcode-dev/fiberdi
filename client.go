package fiberdi

import (
	"github.com/bytedance/sonic"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
)

func New(module IModule, configs ...fiber.Config) *fiber.App {
	config := fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	}

	if len(configs) == 1 {
		config = configs[0]

		config.DisableStartupMessage = true
		config.JSONEncoder = sonic.Marshal
		config.JSONDecoder = sonic.Unmarshal
	}

	if len(configs) > 1 {
		log.Fatalf("it's not supported more than one configuration in fiber")
	}

	app := fiber.New(config)

	module.setInjectableCurrentInAppModuleFn(module.(*Module).setInjectableCurrent)
	module.setInjectablesNames(module.mappedInjectables(module.(*Module).Modules, []string{}))
	module.setExportsNames(module.mappedExports(module.(*Module).Modules, []string{}))

	app = module.start(app)

	return app
}
