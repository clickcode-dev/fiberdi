package cat

import (
	"github.com/gofiber/fiber/v2"
)

type CatController struct {
	CatService *CatService
}

func (controller *CatController) Routes(app *fiber.App) *fiber.App {
	app.Get("/cat", controller.findCats)

	return app
}

func (controller *CatController) findCats(ctx *fiber.Ctx) error {
	return ctx.JSON(controller.CatService.HelloWorld())
}
