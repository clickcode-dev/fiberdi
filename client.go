package fiberdi

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

//	appModule := &fiberdi.Module{}
//
// New creates a new Fiber named instance.
//
//	app := fiberdi.New(appModule)
//
// You can pass optional configuration options by passing a Config struct:
//
//	app := fiber.New(appModule, fiber.Config{
//	    Prefork: true,
//	    ServerHeader: "Fiber",
//	})
//
// ATTENTION: DisableStartupMessage is true and cannot be changed
func New(module *Module, configs ...fiber.Config) *fiber.App {
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

	app = module.start(app)

	return app
}
