package dog

import (
	"github.com/gofiber/fiber/v2"
)

type DogController struct {
	DogService *DogService
}

func (controller DogController) Routes(app *fiber.App) *fiber.App {
	app.Get("/dog", func(c *fiber.Ctx) error {
		return c.JSON(controller.DogService.HelloWorld())
	})

	return app
}
