package cat

import (
	"github.com/gofiber/fiber/v2"
)

type CatController struct {
	CatService *CatService
}

func (controller *CatController) Routes(app *fiber.App) *fiber.App {
	app.Get("/cat", func(c *fiber.Ctx) error {
		return c.JSON(controller.CatService.HelloWorld())
	})

	return app
}
