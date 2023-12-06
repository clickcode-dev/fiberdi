package fiberdi

import (
	"github.com/gofiber/fiber/v2"
)

type IController interface {
	// Routes is needed to create routes automatically
	//
	//	type YourController struct {}
	//
	// 	func (controller YourController) Routes(app *fiber.App) *fiber.App {
	//  	app.Get("/", DoSomething)
	//		return app
	//	}
	Routes(*fiber.App) *fiber.App
}
