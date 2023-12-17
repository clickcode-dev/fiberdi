package elephant

import (
	"github.com/clickcode-dev/fiberdi/.examples/dog"
	"github.com/gofiber/fiber/v2"
)

type ElephantController struct {
	ElephantService *ElephantService
	Service         *dog.DogService
}

func (controller *ElephantController) Routes(app *fiber.App) *fiber.App {
	app.Get("/elephant", func(c *fiber.Ctx) error {
		return c.JSON(controller.ElephantService.HelloWorld())
	})

	return app
}
