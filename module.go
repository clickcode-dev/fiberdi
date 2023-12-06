package fiberdi

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v2"

	"github.com/charmbracelet/log"

	"github.com/sarulabs/di"
)

type IModule interface {
	start(*fiber.App) *fiber.App
	verifyIfIsAnAttemptToInjectController(interface{}, interface{})
	buildDependency(*fiber.App, interface{}, *di.Builder) *fiber.App
	generateDependecies(*fiber.App) *fiber.App
}

// Module manage your dependencies and routes
//
//	appModule := &fiberdi.Module{}
type Module struct {
	// Controller needs function Routes
	//
	//	type YourController struct {}
	//
	// 	func (controller YourController) Routes(app *fiber.App) *fiber.App {
	//  	app.Get("/", DoSomething)
	//		return app
	//	}
	//
	//	appModule := &fiberdi.Module{
	//		Controllers: []fiberdi.IController{
	//			&YourController{},
	//		},
	//	}
	Controllers []IController

	// Injectables persist all your dependencies
	//
	//	type YourService struct {}
	//
	// 	func (controller YourService) DoSomething(ctx *fiber.Ctx) error {
	//  	return ctx.JSON("OK")
	//	}
	//
	//	type YourController struct {
	//		YourService *YourService
	//	}
	//
	// 	func (controller YourController) Routes(app *fiber.App) *fiber.App {
	//  	app.Get("/", controller.YourService.DoSomething)
	//		return app
	//	}
	//
	//	appModule := &fiberdi.Module{
	//		Controllers: []fiberdi.IController{
	//			&YourController{},
	//		},
	//		Injectables: []interface{}{
	//			&YourService{},
	//		},
	//	}
	Injectables []interface{}

	// You can create submodules
	//
	//	type CatService struct {}
	//
	// 	func (controller CatService) DoSomething(ctx *fiber.Ctx) error {
	//  	return ctx.JSON("Cat")
	//	}
	//
	//	type CatController struct {
	//		CatService *CatService
	//	}
	//
	// 	func (controller CatController) Routes(app *fiber.App) *fiber.App {
	//  	app.Get("/", controller.CatService.DoSomething)
	//		return app
	//	}
	//
	//	catModule := &fiberdi.Module{
	//		Controllers: []fiberdi.IController{
	//			&CatController{},
	//		},
	//		Injectables: []interface{}{
	//			&CatService{},
	//		},
	//	}
	//
	//	appModule := &fiberdi.Module{
	//		Modules: []fiberdi.IModule{
	//			catModule,
	//		},
	//	}
	Modules []IModule
}

var loggerGeneral = log.NewWithOptions(os.Stdout, log.Options{
	ReportTimestamp: true,
	Level:           ternary(os.Getenv("GO_ENV") != "production", log.DebugLevel, log.InfoLevel),
})

func (m Module) start(app *fiber.App) *fiber.App {
	for _, module := range m.Modules {
		app = module.start(app)
	}

	return m.generateDependecies(app)
}

func (m Module) verifyIfIsAnAttemptToInjectController(structt interface{}, inject interface{}) {
	_, structtIsController := structt.(IController)
	_, injectIsController := inject.(IController)

	if injectIsController && structtIsController {
		log.Fatalf("you cannot inject %T", inject)
	}
}

func (m Module) injectThirdDependencies(builder *di.Builder) {
	builder.Add(di.Def{
		Name:  "Logger",
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			return newLogger(), nil
		},
	})
}

func (m Module) buildDependency(app *fiber.App, structt interface{}, builder *di.Builder) *fiber.App {
	loggerGeneral.SetPrefix(fmt.Sprintf("%T", structt))

	valueOf := reflect.ValueOf(structt)

	if valueOf.Kind() != reflect.Ptr {
		log.Fatalf("are you sure? maybe %T is not a pointer, please check", structt)
	}

	pointer := reflect.Indirect(valueOf)

	builder.Add(di.Def{
		Name:  pointer.Type().Name(),
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			for i := 0; i < pointer.NumField(); i++ {
				field := pointer.Type().Field(i)
				name := field.Name

				tagDi := getStructTag(field, "di")

				if tagDi == "ignore" {
					loggerGeneral.Debugf("Ignoring field %s because tag ignore is set", name)
					continue
				}

				inject, err := ctn.SafeGet(name)

				if err != nil {
					log.Fatalf("tip: are you trying to access %T inside %T before inject in module?\n\n%v", inject, structt, err)
				}

				if logger, ok := inject.(*log.Logger); ok {
					logger.SetPrefix(pointer.Type().Name())
				}

				injectValueOf := reflect.ValueOf(inject)

				m.verifyIfIsAnAttemptToInjectController(structt, inject)

				valueOf.Elem().Field(i).Set(injectValueOf)
			}

			if controller, ok := structt.(IController); ok {
				app = controller.Routes(app)
			}

			loggerGeneral.Debugf("dependencies initialized")

			return structt, nil
		},
	})

	return app
}

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func (m Module) generateDependecies(app *fiber.App) *fiber.App {
	builder, err := di.NewBuilder(di.Request)

	if err != nil {
		panic(err)
	}

	m.injectThirdDependencies(builder)

	for _, injectable := range m.Injectables {
		app = m.buildDependency(app, injectable, builder)
	}

	for _, controller := range m.Controllers {
		app = m.buildDependency(app, controller, builder)
	}

	build := builder.Build()

	for _, controller := range m.Controllers {
		typee := reflect.Indirect(reflect.ValueOf(controller)).Type()

		build.Get(typee.Name())
	}

	return app
}
