package fiberdi

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type TestServiceNotMapped struct{}

type TestController struct {
	Service *TestService
}

type TestOtherController struct {
	Service *TestService
}

type TestControllerWithInjectableNonPointer struct {
	Service TestService
}

type TestControllerWithInjectableNotMapped struct {
	Service TestServiceNotMapped
}

type TestControllerWithIgnoreField struct {
	Name string `di:"ignore"`
}

type TestService struct {
	Log *log.Logger
}

func (controller *TestController) Routes(app *fiber.App) *fiber.App {
	return app
}

func (controller *TestOtherController) Routes(app *fiber.App) *fiber.App {
	return app
}

func (controller *TestControllerWithIgnoreField) Routes(app *fiber.App) *fiber.App {
	return app
}

func (controller *TestControllerWithInjectableNonPointer) Routes(app *fiber.App) *fiber.App {
	return app
}

func (controller *TestControllerWithInjectableNotMapped) Routes(app *fiber.App) *fiber.App {
	return app
}

func TestModuleWithoutChild(t *testing.T) {
	controller := new(TestController)
	service := new(TestService)

	module := &Module{
		Controllers: []IController{
			controller,
		},
		Injectables: []interface{}{
			service,
		},
	}

	_ = New(module)

	assert.NotEqual(t, nil, controller.Service)
}

func TestModuleWithIgnoreField(t *testing.T) {
	controller := new(TestControllerWithIgnoreField)

	module := &Module{
		Controllers: []IController{
			controller,
		},
	}

	_ = New(module)

	assert.Equal(t, "", controller.Name)
}

func TestModulePanicWhenInjectableIsNotPointer(t *testing.T) {
	assert.Panics(t, func() {
		controller := new(TestControllerWithInjectableNonPointer)

		module := &Module{
			Controllers: []IController{
				controller,
			},
		}
		New(module)
	})
}

func TestModulePanicWhenInjectableNotMapped(t *testing.T) {
	assert.Panics(t, func() {
		controller := new(TestControllerWithInjectableNotMapped)

		module := &Module{
			Controllers: []IController{
				controller,
			},
		}
		New(module)
	})
}

func TestModuleWithChild(t *testing.T) {
	controller := new(TestController)
	service := new(TestService)

	childModule := &Module{
		Controllers: []IController{
			controller,
		},
		Injectables: []interface{}{
			service,
		},
	}

	module := &Module{
		Modules: []IModule{
			childModule,
		},
	}

	_ = New(module)

	assert.NotEqual(t, nil, controller.Service)
}

func TestModuleWithImportModule(t *testing.T) {
	controller := new(TestController)
	otherController := new(TestOtherController)

	service := new(TestService)

	childModule := &Module{
		Controllers: []IController{
			controller,
		},
		Injectables: []interface{}{
			service,
		},
		Exports: []interface{}{
			service,
		},
	}

	otherChildModule := &Module{
		Controllers: []IController{
			otherController,
		},
		Injectables: []interface{}{
			service,
		},
		Imports: []IModule{
			childModule,
		},
	}

	module := &Module{
		Modules: []IModule{
			childModule,
			otherChildModule,
		},
	}

	_ = New(module)

	assert.NotEqual(t, nil, controller.Service)
	assert.NotEqual(t, nil, otherController.Service)
}
