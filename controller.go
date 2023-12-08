package fiberdi

import (
	"github.com/gofiber/fiber/v2"
)

type IController interface {
	Routes(*fiber.App) *fiber.App
}
