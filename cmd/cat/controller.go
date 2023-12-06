package cat

import (
	"github.com/clickcode-dev/fiberdi/cmd/dog"
	"github.com/gofiber/fiber/v2"
)

type CatController struct {
	CatService *CatService
	DogService *dog.DogService
}

func (controller CatController) Routes(app *fiber.App) *fiber.App {
	app.Get("/cat", func(c *fiber.Ctx) error {
		return c.JSON(controller.DogService.HelloWorld())
	})

	return app
}
